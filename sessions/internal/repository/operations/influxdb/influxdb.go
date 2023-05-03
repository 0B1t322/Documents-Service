package influxdb

import (
	"context"
	"fmt"
	"github.com/0B1t322/Documents-Service/sessions/internal/core/models"
	"github.com/0B1t322/Documents-Service/sessions/internal/core/operations/elements"
	repository "github.com/0B1t322/Documents-Service/sessions/internal/repository/operations"
	"github.com/go-faster/errors"
	"github.com/google/uuid"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/query"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/samber/lo"
	"sort"
	"strconv"
	"time"
)

type Repository struct {
	client influxdb2.Client
	Org    string
	Bucket string
}

func New(client influxdb2.Client, org string, bucket string) *Repository {
	return &Repository{client: client, Org: org, Bucket: bucket}
}

func (r Repository) SaveRevision(
	ctx context.Context,
	documentId uuid.UUID,
	operations []models.RevisionOperation,
) (int, error) {
	//writeApi := r.client.WriteAPI(r.Org, r.Bucket)

	revisionId, err := r.GetLastRevision(ctx, documentId)
	if err != nil {
		return 0, err
	}

	// increment
	revisionId += 1

	writeApiBlocking := r.client.WriteAPIBlocking(r.Org, r.Bucket)

	writeApiBlocking.EnableBatching()

	points := operationToPoints(documentId, revisionId, operations)

	if err := writeApiBlocking.WritePoint(ctx, points...); err != nil {
		return 0, err
	}

	if err := writeApiBlocking.Flush(ctx); err != nil {
		return 0, err
	}

	return revisionId, nil
}

func operationToPoints(
	documentId uuid.UUID, revisionId int,
	operations []models.RevisionOperation,
) (points []*write.Point) {
	for _, op := range operations {
		p := operationToPoint(documentId, revisionId, op)
		points = append(points, p)
	}

	return
}

func operationToPoint(
	documentId uuid.UUID, revisionId int,
	operation models.RevisionOperation,
) *write.Point {
	measurement := "revisions"

	tags := map[string]string{
		"documentId": documentId.String(),
		"revisionId": fmt.Sprint(revisionId),
	}

	fields := map[string]interface{}{}
	t := operation.Type()
	fields["type"] = t

	switch t {
	case models.RevisionOperationTypeInsert:
		fields["structuralElementIndex"] = operation.Insert.StructuralElementIndex
		fields["paragraphElementIndex"] = operation.Insert.InsertParagraphElement.Index
		fields["localIndex"] = operation.Insert.InsertText.Index
		fields["content"] = operation.Insert.InsertText.Text
	case models.RevisionOperationTypeDelete:
		fields["structuralElementIndex"] = operation.Delete.StructuralElementIndex
		fields["paragraphElementIndex"] = operation.Delete.DeleteParagraphElement.Index
		fields["localIndex"] = operation.Delete.DeleteText.Index
		fields["content"] = operation.Delete.DeleteText.Text
	}

	p := influxdb2.NewPoint(measurement, tags, fields, time.Now())

	return p
}

func (r Repository) GetLastRevision(
	ctx context.Context,
	documentId uuid.UUID,
) (int, error) {
	queryApi := r.client.QueryAPI(r.Org)

	res, err := queryApi.Query(
		ctx,
		fmt.Sprintf(
			`
		import "influxdata/influxdb/schema"
				
		schema.tagValues(
			bucket: "%s", 
			tag: "revisionId", 
			predicate: (r) => 
				r.documentId == "%s",
			)
			|> last()`,
			r.Bucket,
			documentId,
		),
	)
	if err != nil {
		return 0, err
	}
	defer res.Close()

	if !res.Next() {
		return -1, nil
	}

	revision := res.Record().Value().(string)

	revisionId, _ := strconv.Atoi(revision)

	return revisionId, nil
}

func (r Repository) GetRevisions(ctx context.Context, documentId uuid.UUID) ([]models.DocumentRevision, error) {
	queryApi := r.client.QueryAPI(r.Org)

	res, err := queryApi.Query(
		ctx,
		fmt.Sprintf(
			`
				from(bucket: "%s")
					|> range(start: 0)
					|> filter(fn: (r) => 
						r.documentId == "%s",
					)
					|> keep(columns: ["_field", "_value", "documentId", "revisionId", "_time"])
					|> pivot(
						columnKey: ["_field"], 
						rowKey: ["documentId","revisionId","_time"], 
						valueColumn: "_value"
					)
					|> drop(columns: ["_time"])`,
			r.Bucket,
			documentId,
		),
	)
	if err != nil {
		return []models.DocumentRevision{}, err
	}

	revs := scanDocumentRevisions(res)

	if len(revs) == 0 {
		return revs, repository.ErrNotFoundRevisionsForDocument
	}

	return revs, nil
}

func (r Repository) GetRevisionsAfter(
	ctx context.Context,
	documentId uuid.UUID,
	revId int,
) ([]models.DocumentRevision, error) {
	queryApi := r.client.QueryAPI(r.Org)

	res, err := queryApi.Query(
		ctx,
		fmt.Sprintf(
			`
				from(bucket: "%s")
					|> range(start: 0)
					|> filter(fn: (r) => 
						r.documentId == "%s" and
						r.revisionId > "%v",
					)
					|> keep(columns: ["_field", "_value", "documentId", "revisionId", "_time"])
					|> pivot(
						columnKey: ["_field"], 
						rowKey: ["documentId","revisionId","_time"], 
						valueColumn: "_value"
					)
					|> drop(columns: ["_time"])`,
			r.Bucket,
			documentId,
			revId,
		),
	)
	if err != nil {
		return []models.DocumentRevision{}, err
	}

	revs := scanDocumentRevisions(res)

	if len(revs) == 0 {
		return revs, repository.ErrNotFoundRevisionsForDocument
	}

	return revs, nil
}

func scanDocumentRevisions(r *api.QueryTableResult) (revs []models.DocumentRevision) {
	m := map[string]models.DocumentRevision{}
	for r.Next() {
		rev := scanDocumentRevision(r.Record())
		mapKey := fmt.Sprintf("%s:%v", rev.DocumentID, rev.RevisionID)

		if _, find := m[mapKey]; !find {
			m[mapKey] = rev
			continue
		}
		value := m[mapKey]
		value.Operations = append(value.Operations, rev.Operations...)
		m[mapKey] = value
	}

	for _, v := range m {
		revs = append(revs, v)
	}

	sort.Slice(
		revs, func(i, j int) bool {
			return revs[i].RevisionID < revs[j].RevisionID
		},
	)

	return
}

func scanDocumentRevision(r *query.FluxRecord) models.DocumentRevision {
	var (
		documentId             = uuid.MustParse(r.ValueByKey("documentId").(string))
		revisionId, _          = strconv.Atoi(r.ValueByKey("revisionId").(string))
		operationType          = models.RevisionOperationType(r.ValueByKey("type").(string))
		structuralElementIndex = int(r.ValueByKey("structuralElementIndex").(int64))
		paragraphElementIndex  = int(r.ValueByKey("paragraphElementIndex").(int64))
		index                  = int(r.ValueByKey("localIndex").(int64))
		content                = r.ValueByKey("content").(string)
	)

	var operation models.RevisionOperation
	{
		switch operationType {
		case models.RevisionOperationTypeInsert:
			operation.Insert = lo.ToPtr(
				elements.NewInsertStructuralElement(
					structuralElementIndex, paragraphElementIndex,
					index, content,
				),
			)
		case models.RevisionOperationTypeDelete:
			operation.Delete = lo.ToPtr(
				elements.NewDeleteStructuralElement(
					structuralElementIndex, paragraphElementIndex,
					index, content,
				),
			)
		}
	}

	return models.DocumentRevision{
		DocumentID: documentId,
		RevisionID: revisionId,
		Operations: []models.RevisionOperation{operation},
	}
}

func (r Repository) IsNotFound(err error) bool {
	return errors.Is(err, repository.ErrNotFoundRevisionsForDocument)
}
