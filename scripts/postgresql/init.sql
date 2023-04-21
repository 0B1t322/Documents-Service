create sequence "StructuralElements_id_seq";

alter sequence "StructuralElements_id_seq" owner to postgres;

create sequence "ParagraphStyle_Id_seq";

alter sequence "ParagraphStyle_Id_seq" owner to postgres;

create sequence "ParagraphStyle_Id_seq1";

alter sequence "ParagraphStyle_Id_seq1" owner to postgres;

create type baseline_offset as enum ('BASELINE_OFFSET_UNSPECIFIED', 'NONE', 'SUPERSCRIPT', 'SUBSCRIPT');

alter type baseline_offset owner to postgres;

create type unit as enum ('UNIT_UNSPECIFIED', 'PT');

alter type unit owner to postgres;

create type dashstyle as enum ('DASH_STYLE_UNSPECIFIED', 'SOLID', 'DOT', 'DASH');

alter type dashstyle owner to postgres;

create type indexes as
(
    "StartIndex" bigint,
    "EndIndex"   bigint
);

alter type indexes owner to postgres;

create type color as
(
    "Red"   double precision,
    "Blue"  double precision,
    "Green" double precision
);

alter type color owner to postgres;

create type dimensions as
(
    "Magnitude" double precision,
    "Unit"      unit
);

alter type dimensions owner to postgres;

create type embedded_object_border as
(
    "Color"     color,
    "Width"     dimensions,
    "DashStyle" dashstyle
);

alter type embedded_object_border owner to postgres;

create type size as
(
    "Height" dimensions,
    "Width"  dimensions
);

alter type size owner to postgres;

create type crop_properties as
(
    "OffsetLeft"   double precision,
    "OffsetRight"  double precision,
    "OffsetTop"    double precision,
    "OffsetBottom" double precision,
    "Angle"        double precision
);

comment on type crop_properties is 'The crop properties of an image. Ange in radians. If all offsets and rotation are 0, the image is not cropped. If the offset is in the interval (0, 1), the corresponding edge of crop rectangle is positioned inside of the image''s original bounding rectangle. If the offset is negative or greater than 1, the corresponding edge of crop rectangle is positioned outside of the image''s original bounding rectangle.';

alter type crop_properties owner to postgres;

create type section_type as enum ('UNSPECIFIED', 'CONTINUOUS', 'NEXT_PAGE');

alter type section_type owner to postgres;

create type column_separator_style as enum ('NONE', 'BETWEEN_EACH_COLUMN');

alter type column_separator_style owner to postgres;

create type content_direction as enum ('CONTENT_DIRECTION_UNSPECIFIED', 'LEFT_TO_RIGHT', 'RIGHT_TO_LEFT');

alter type content_direction owner to postgres;

create type section_column_properties as
(
    "Width"      dimensions,
    "PaddingEnd" dimensions
);

alter type section_column_properties owner to postgres;

create type paragraph_border as
(
    "Color"     color,
    "Width"     dimensions,
    "Padding"   dimensions,
    "DashStyle" dashstyle
);

alter type paragraph_border owner to postgres;

create type named_style_type as enum ('NORMAL_TEXT', 'TITLE', 'SUBTITLE', 'HEADING_1', 'HEADING_2', 'HEADING_3', 'HEADING_4', 'HEADING_5', 'HEADING_6');

alter type named_style_type owner to postgres;

create type alignment as enum ('UNSPECIFIED', 'START', 'CENTER', 'END', 'JUSTIFIED');

alter type alignment owner to postgres;

create type spacing_mode as enum ('UNSPECIFIED', 'NEVER_COLLAPSE', 'COLLAPSE_LISTS');

alter type spacing_mode owner to postgres;

create type tab_stop_alignment as enum ('UNSPECIFIED', 'START', 'CENTER', 'END');

alter type tab_stop_alignment owner to postgres;

create type tab_stop as
(
    "Offset"    dimensions,
    "Alignment" tab_stop_alignment
);

alter type tab_stop owner to postgres;

create table "Body"
(
    "Id" uuid default gen_random_uuid() not null
        constraint "Body_pk"
            primary key
);

alter table "Body"
    owner to postgres;

create table "Link"
(
    "Id"         bigint generated always as identity
        constraint "Link_pk"
            primary key,
    "Url"        text,
    "BookmarkId" uuid,
    "HeadingId"  uuid,
    constraint check_that_only_one_field_not_null
        check (num_nonnulls("Url", "BookmarkId", "HeadingId") = 1)
);

alter table "Link"
    owner to postgres;

create table "TextStyle"
(
    "Id"              bigint                                                                 not null
        constraint "TextStyle_pk"
            primary key,
    "Bold"            boolean,
    "Italic"          boolean,
    "Underline"       boolean,
    "Strikethrough"   boolean,
    "FontFamily"      text,
    "FontWeight"      integer
        constraint "FontWeightBetween_100_900"
            check (("FontWeight" >= 100) AND ("FontWeight" <= 900)),
    "LinkId"          bigint
        constraint "TextStyle_Link_Id_fk"
            references "Link",
    "BaselineOffset"  baseline_offset default 'BASELINE_OFFSET_UNSPECIFIED'::baseline_offset not null,
    "BackgroundColor" color,
    "ForegroundColor" color,
    "FontSize"        dimensions,
    "SmallCaps"       boolean
);

alter table "TextStyle"
    owner to postgres;

create table "TextRun"
(
    "Id"          bigint generated always as identity
        constraint "TextRun_pk"
            primary key,
    "Content"     text default '\n'::text not null,
    "TextStyleId" bigint
        constraint "TextRun_TextStyle_Id_fk"
            references "TextStyle"
);

alter table "TextRun"
    owner to postgres;

create table "InlineObjectsElements"
(
    "Id"             bigint generated always as identity
        constraint "InlineObjectsElements_pk"
            primary key,
    "InlineObjectId" uuid not null,
    "TextStyleId"    bigint
        constraint "InlineObjectsElements_TextStyle_Id_fk"
            references "TextStyle"
);

alter table "InlineObjectsElements"
    owner to postgres;

create table "ImageProperties"
(
    "Id"             bigint generated always as identity
        constraint "ImageProperties_pk"
            primary key,
    "SourceUri"      text,
    "Brightness"     double precision default 0                                                                                                                       not null,
    "Contrast"       double precision default 0                                                                                                                       not null,
    "Transparency"   double precision default 0                                                                                                                       not null,
    "Angle"          double precision default 0                                                                                                                       not null,
    "CropProperties" crop_properties  default ROW ((0)::double precision, (0)::double precision, (0)::double precision, (0)::double precision, (0)::double precision) not null,
    "ContentUri"     text                                                                                                                                             not null
);

comment on column "ImageProperties"."Angle" is 'In radians';

alter table "ImageProperties"
    owner to postgres;

create table "EmbeddedObjects"
(
    "Id"                   bigint generated always as identity
        constraint "EmbeddedObjects_pk"
            primary key,
    "Title"                text                   default ''::text                                                not null,
    "Description"          text                   default ''::text                                                not null,
    "EmbeddedObjectBorder" embedded_object_border default ROW (NULL::color, NULL::dimensions, 'SOLID'::dashstyle) not null,
    "Size"                 size                                                                                   not null,
    "MarginTop"            dimensions             default ROW ((0)::double precision, 'PT'::unit),
    "MarginBottom"         dimensions             default ROW ((0)::double precision, 'PT'::unit),
    "MarginRight"          dimensions             default ROW ((0)::double precision, 'PT'::unit),
    "MarginLeft"           dimensions             default ROW ((0)::double precision, 'PT'::unit),
    "ImagePropertiesId"    bigint                                                                                 not null
        constraint "EmbeddedObjects_ImageProperties_Id_fk"
            references "ImageProperties"
);

alter table "EmbeddedObjects"
    owner to postgres;

create table "InlineObjectProperties"
(
    "Id"               bigint generated always as identity
        constraint "InlineObjectProperties_pk"
            primary key,
    "EmbeddedObjectId" bigint not null
        constraint "InlineObjectProperties___fk"
            references "EmbeddedObjects"
);

alter table "InlineObjectProperties"
    owner to postgres;

create table "PageBreak"
(
    "Id"          bigint generated always as identity
        constraint "PageBreak_pk"
            primary key,
    "TextStyleId" bigint
        constraint "PageBreak_TextStyle_Id_fk"
            references "TextStyle"
);

alter table "PageBreak"
    owner to postgres;

create table "SectionStyle"
(
    "Id"                       bigint generated always as identity
        constraint "SectionStyle_pk"
            primary key,
    "ColumnSeparatorStyle"     column_separator_style default 'NONE'::column_separator_style     not null,
    "ContentDirection"         content_direction      default 'LEFT_TO_RIGHT'::content_direction not null,
    "MarginTop"                dimensions,
    "MarginLeft"               dimensions,
    "MarginBottom"             dimensions,
    "MarginRight"              dimensions,
    "MarginHeader"             dimensions,
    "MarginFooter"             dimensions,
    "SectionType"              section_type           default 'CONTINUOUS'::section_type         not null,
    "DefaultHeaderId"          uuid,
    "DefaultFooterId"          uuid,
    "FirstPageHeaderId"        uuid,
    "FirstPageFooterId"        uuid,
    "EvenPageHeaderId"         uuid,
    "EvenPageFooterId"         uuid,
    "UseFirstPageHeaderFooter" boolean                default false,
    "PageNumberStart"          bigint,
    "ColumnProperties"         section_column_properties[]
);

alter table "SectionStyle"
    owner to postgres;

create table "SectionBreak"
(
    "Id"             bigint generated always as identity
        constraint "SectionBreak_pk"
            primary key,
    "SectionStyleId" bigint not null
        constraint "SectionBreak_SectionStyle_Id_fk"
            references "SectionStyle"
);

alter table "SectionBreak"
    owner to postgres;

create table "Equation"
(
    "Id"          bigint generated always as identity
        constraint "Equation_pk"
            primary key,
    "Content"     text,
    "TextStyleId" bigint
        constraint "Equation_TextStyle_Id_fk"
            references "TextStyle"
);

alter table "Equation"
    owner to postgres;

create table "DocumentStyles"
(
    "Id"                           bigint generated always as identity
        constraint "DocumentStyles_pk"
            primary key,
    "DefaultHeaderId"              uuid,
    "DefaultFooterId"              uuid,
    "EvenPageHeaderId"             uuid,
    "EventPageFooterId"            uuid,
    "FirstPageHeaderId"            uuid,
    "FirstPageFooterId"            uuid,
    "UseFirstPageHeaderFooter"     boolean    default false                                    not null,
    "UseEvenPageHeaderFooter"      boolean    default false,
    "PageNumberStart"              bigint     default 1                                        not null,
    "MarginTop"                    dimensions default ROW ((72)::double precision, 'PT'::unit) not null,
    "MarginBottom"                 dimensions default ROW ((72)::double precision, 'PT'::unit) not null,
    "MarginRight"                  dimensions default ROW ((72)::double precision, 'PT'::unit) not null,
    "MarginLeft"                   dimensions default ROW ((72)::double precision, 'PT'::unit) not null,
    "PageSize"                     size                                                        not null,
    "MarginHeader"                 dimensions                                                  not null,
    "UseCustomHeaderFooterMargins" boolean    default false                                    not null
);

alter table "DocumentStyles"
    owner to postgres;

alter sequence "ParagraphStyle_Id_seq" owned by "DocumentStyles"."Id";

create table "Documents"
(
    "Id"              uuid default gen_random_uuid() not null
        constraint "Documents_pk"
            primary key,
    "Title"           text,
    "BodyId"          uuid                           not null
        constraint "Documents_Body_Id_fk"
            references "Body"
            on delete cascade,
    "DocumentStyleId" bigint                         not null
        constraint "Documents_DocumentStyles_Id_fk"
            references "DocumentStyles"
            on delete cascade
);

alter table "Documents"
    owner to postgres;

create table "InlineObjects"
(
    "Id"                       bigint not null
        constraint "InlineObjects_pk"
            primary key,
    "ObjectId"                 uuid   not null,
    "InlineObjectPropertiesId" bigint not null
        constraint "InlineObjects_InlineObjectProperties_Id_fk"
            references "InlineObjectProperties",
    "DocumentId"               uuid   not null
        constraint "InlineObjects_Documents_Id_fk"
            references "Documents"
);

alter table "InlineObjects"
    owner to postgres;

create table "ParagraphStyle"
(
    "Id"                  bigint generated always as identity
        constraint "ParagraphStyle_pk"
            primary key,
    "HeadingId"           uuid,
    "Aligment"            alignment         default 'UNSPECIFIED'::alignment                                                                                                not null,
    "LineSpacing"         integer           default 100                                                                                                                     not null
        constraint "LineSpacing_Between_6_10000"
            check (("LineSpacing" >= 6) AND ("LineSpacing" <= 10000)),
    "Direction"           content_direction default 'LEFT_TO_RIGHT'::content_direction                                                                                      not null,
    "SpacingMode"         spacing_mode      default 'UNSPECIFIED'::spacing_mode                                                                                             not null,
    "SpaceAbove"          dimensions        default ROW ((0)::double precision, 'PT'::unit)                                                                                 not null,
    "SpaceBelow"          dimensions        default ROW ((0)::double precision, 'PT'::unit)                                                                                 not null,
    "BorderBetween"       paragraph_border  default ROW (NULL::color, ROW ((0)::double precision, 'PT'::unit), ROW ((0)::double precision, 'PT'::unit), 'SOLID'::dashstyle) not null,
    "BorderTop"           paragraph_border  default ROW (NULL::color, ROW ((0)::double precision, 'PT'::unit), ROW ((0)::double precision, 'PT'::unit), 'SOLID'::dashstyle) not null,
    "BorderBottom"        paragraph_border  default ROW (NULL::color, ROW ((0)::double precision, 'PT'::unit), ROW ((0)::double precision, 'PT'::unit), 'SOLID'::dashstyle) not null,
    "BorderLeft"          paragraph_border  default ROW (NULL::color, ROW ((0)::double precision, 'PT'::unit), ROW ((0)::double precision, 'PT'::unit), 'SOLID'::dashstyle) not null,
    "BorderRight"         paragraph_border  default ROW (NULL::color, ROW ((0)::double precision, 'PT'::unit), ROW ((0)::double precision, 'PT'::unit), 'SOLID'::dashstyle) not null,
    "IndentFirstLine"     dimensions,
    "IndentStart"         dimensions,
    "IndentEnd"           paragraph_border,
    "TabStops"            tab_stop[]        default ARRAY []::tab_stop[],
    "KeepLinesTogether"   boolean,
    "KeepWithNext"        boolean,
    "AvoidWidowAndOrphan" boolean,
    "Shading"             color,
    "PageBreakBefore"     boolean
);

alter table "ParagraphStyle"
    owner to postgres;

alter sequence "ParagraphStyle_Id_seq1" owned by "ParagraphStyle"."Id";

create table "Paragraphs"
(
    "Id"               bigint generated always as identity
        constraint "Paragraphs_pk"
            primary key,
    "ParagraphStyleId" bigint not null
        constraint "Paragraphs_ParagraphStyle_Id_fk"
            references "ParagraphStyle"
);

alter table "Paragraphs"
    owner to postgres;

create table "StructuralElements"
(
    "Id"             bigint generated always as identity
        constraint "StructuralElements_pk"
            primary key,
    "BodyId"         uuid                                           not null
        constraint "StructuralElements_Body_Id_fk"
            references "Body",
    "ParagraphId"    bigint
        constraint "StructuralElements_Paragraphs_Id_fk"
            references "Paragraphs",
    "Indexes"        indexes default ROW ((0)::bigint, (0)::bigint) not null,
    "SectionBreakId" bigint
        constraint "StructuralElements_SectionBreak_Id_fk"
            references "SectionBreak",
    constraint only_one_element_not_null
        check (num_nonnulls("ParagraphId", "SectionBreakId") = 1)
);

alter table "StructuralElements"
    owner to postgres;

alter sequence "StructuralElements_id_seq" owned by "StructuralElements"."Id";

create unique index "StructuralElements_BodyId_Indexes_uindex"
    on "StructuralElements" ("BodyId", "Indexes");

create table "ParagraphElements"
(
    "Id"                    bigint generated always as identity
        constraint "ParagraphElements_pk"
            primary key,
    "ParagraphId"           bigint  not null
        constraint "ParagraphElements_Paragraphs_Id_fk"
            references "Paragraphs",
    "TextRunId"             bigint
        constraint "ParagraphElements_TextRun_Id_fk"
            references "TextRun",
    "InlineObjectElementId" bigint
        constraint "ParagraphElements_InlineObjectsElements_Id_fk"
            references "InlineObjectsElements",
    "Indexes"               indexes not null,
    "PageBreakId"           bigint
        constraint "ParagraphElements_PageBreak_Id_fk"
            references "PageBreak",
    "EquationId"            bigint
        constraint "ParagraphElements_Equation_Id_fk"
            references "Equation",
    constraint only_one_element_not_null
        check (num_nonnulls("TextRunId", "InlineObjectElementId", "PageBreakId", "EquationId") = 1)
);

alter table "ParagraphElements"
    owner to postgres;

create unique index "ParagraphElements_Id_Indexes_uindex"
    on "ParagraphElements" ("Id", "Indexes");

create table "Styles"
(
    "Id"               uuid default gen_random_uuid() not null
        constraint "Styles_pk"
            primary key,
    "Name"             text                           not null,
    "ParagraphStyleId" bigint                         not null
        constraint "Styles_ParagraphStyle_Id_fk"
            references "ParagraphStyle",
    "TextStyleId"      bigint                         not null
        constraint "Styles_TextStyle_Id_fk"
            references "TextStyle"
);

alter table "Styles"
    owner to postgres;

create unique index "Styles_Name_uindex"
    on "Styles" ("Name");

create unique index "ParagraphStyle_HeadingId_uindex"
    on "ParagraphStyle" ("HeadingId");

create table "DocumentsStyles"
(
    "Id"         bigint generated always as identity
        constraint "DocumentsStyles_pk"
            primary key,
    "DocumentId" uuid not null
        constraint "DocumentsStyles_Documents_Id_fk"
            references "Documents",
    "StyleId"    uuid not null
        constraint "DocumentsStyles_Styles_Id_fk"
            references "Styles"
);

alter table "DocumentsStyles"
    owner to postgres;

