create table if not exists public."Body"
(
    "Id" uuid default gen_random_uuid() not null
        constraint "Body_pk"
            primary key
);

alter table public."Body"
    owner to postgres;

create table if not exists public."Link"
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

alter table public."Link"
    owner to postgres;

create table if not exists public."TextStyle"
(
    "Id"              bigint generated always as identity
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
            references public."Link",
    "BaselineOffset"  baseline_offset default 'BASELINE_OFFSET_UNSPECIFIED'::baseline_offset not null,
    "BackgroundColor" color,
    "ForegroundColor" color,
    "FontSize"        dimensions,
    "SmallCaps"       boolean
);

alter table public."TextStyle"
    owner to postgres;

create table if not exists public."TextRun"
(
    "Id"          bigint generated always as identity
        constraint "TextRun_pk"
            primary key,
    "Content"     text default '\n'::text not null,
    "TextStyleId" bigint
        constraint "TextRun_TextStyle_Id_fk"
            references public."TextStyle"
);

alter table public."TextRun"
    owner to postgres;

create table if not exists public."InlineObjectsElements"
(
    "Id"             bigint generated always as identity
        constraint "InlineObjectsElements_pk"
            primary key,
    "InlineObjectId" uuid not null,
    "TextStyleId"    bigint
        constraint "InlineObjectsElements_TextStyle_Id_fk"
            references public."TextStyle"
);

alter table public."InlineObjectsElements"
    owner to postgres;

create table if not exists public."ImageProperties"
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

comment on column public."ImageProperties"."Angle" is 'In radians';

alter table public."ImageProperties"
    owner to postgres;

create table if not exists public."EmbeddedObjects"
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
            references public."ImageProperties"
);

alter table public."EmbeddedObjects"
    owner to postgres;

create table if not exists public."InlineObjectProperties"
(
    "Id"               bigint generated always as identity
        constraint "InlineObjectProperties_pk"
            primary key,
    "EmbeddedObjectId" bigint not null
        constraint "InlineObjectProperties___fk"
            references public."EmbeddedObjects"
);

alter table public."InlineObjectProperties"
    owner to postgres;

create table if not exists public."PageBreak"
(
    "Id"          bigint generated always as identity
        constraint "PageBreak_pk"
            primary key,
    "TextStyleId" bigint
        constraint "PageBreak_TextStyle_Id_fk"
            references public."TextStyle"
);

alter table public."PageBreak"
    owner to postgres;

create table if not exists public."SectionStyle"
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

alter table public."SectionStyle"
    owner to postgres;

create table if not exists public."SectionBreak"
(
    "Id"             bigint generated always as identity
        constraint "SectionBreak_pk"
            primary key,
    "SectionStyleId" bigint
        constraint "SectionBreak_SectionStyle_Id_fk"
            references public."SectionStyle"
);

alter table public."SectionBreak"
    owner to postgres;

create table if not exists public."Equation"
(
    "Id"          bigint generated always as identity
        constraint "Equation_pk"
            primary key,
    "Content"     text,
    "TextStyleId" bigint
        constraint "Equation_TextStyle_Id_fk"
            references public."TextStyle"
);

alter table public."Equation"
    owner to postgres;

create table if not exists public."DocumentStyles"
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
    "MarginHeader"                 dimensions default ROW ((36)::double precision, 'PT'::unit) not null,
    "UseCustomHeaderFooterMargins" boolean    default false                                    not null
);

alter table public."DocumentStyles"
    owner to postgres;

create table if not exists public."Documents"
(
    "Id"              uuid default gen_random_uuid() not null
        constraint "Documents_pk"
            primary key,
    "Title"           text,
    "BodyId"          uuid                           not null
        constraint "Documents_Body_Id_fk"
            references public."Body"
            on delete cascade,
    "DocumentStyleId" bigint                         not null
        constraint "Documents_DocumentStyles_Id_fk"
            references public."DocumentStyles"
            on delete cascade
);

alter table public."Documents"
    owner to postgres;

create table if not exists public."InlineObjects"
(
    "Id"                       bigint not null
        constraint "InlineObjects_pk"
            primary key,
    "ObjectId"                 uuid   not null,
    "InlineObjectPropertiesId" bigint not null
        constraint "InlineObjects_InlineObjectProperties_Id_fk"
            references public."InlineObjectProperties",
    "DocumentId"               uuid   not null
        constraint "InlineObjects_Documents_Id_fk"
            references public."Documents"
);

alter table public."InlineObjects"
    owner to postgres;

create table if not exists public."ParagraphStyle"
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

alter table public."ParagraphStyle"
    owner to postgres;

create unique index if not exists "ParagraphStyle_HeadingId_uindex"
    on public."ParagraphStyle" ("HeadingId");

create table if not exists public."Paragraphs"
(
    "Id"               bigint generated always as identity
        constraint "Paragraphs_pk"
            primary key,
    "ParagraphStyleId" bigint
        constraint "Paragraphs_ParagraphStyle_Id_fk"
            references public."ParagraphStyle"
);

alter table public."Paragraphs"
    owner to postgres;

create table if not exists public."StructuralElements"
(
    "Id"             bigint generated always as identity
        constraint "StructuralElements_pk"
            primary key,
    "BodyId"         uuid    not null
        constraint "StructuralElements_Body_Id_fk"
            references public."Body",
    "ParagraphId"    bigint
        constraint "StructuralElements_Paragraphs_Id_fk"
            references public."Paragraphs",
    "SectionBreakId" bigint
        constraint "StructuralElements_SectionBreak_Id_fk"
            references public."SectionBreak",
    "Index"          integer not null,
    constraint only_one_element_not_null
        check (num_nonnulls("ParagraphId", "SectionBreakId") = 1)
);

alter table public."StructuralElements"
    owner to postgres;

create unique index if not exists "StructuralElements_ParagraphId_uindex"
    on public."StructuralElements" ("ParagraphId");

create table if not exists public."ParagraphElements"
(
    "Id"                    bigint generated always as identity
        constraint "ParagraphElements_pk"
            primary key,
    "ParagraphId"           bigint not null
        constraint "ParagraphElements_Paragraphs_Id_fk"
            references public."Paragraphs",
    "TextRunId"             bigint
        constraint "ParagraphElements_TextRun_Id_fk"
            references public."TextRun",
    "InlineObjectElementId" bigint
        constraint "ParagraphElements_InlineObjectsElements_Id_fk"
            references public."InlineObjectsElements",
    "PageBreakId"           bigint
        constraint "ParagraphElements_PageBreak_Id_fk"
            references public."PageBreak",
    "EquationId"            bigint
        constraint "ParagraphElements_Equation_Id_fk"
            references public."Equation",
    "Index"                 integer,
    constraint only_one_element_not_null
        check (num_nonnulls("TextRunId", "InlineObjectElementId", "PageBreakId", "EquationId") = 1)
);

alter table public."ParagraphElements"
    owner to postgres;

create table if not exists public."Styles"
(
    "Id"               uuid default gen_random_uuid() not null
        constraint "Styles_pk"
            primary key,
    "Name"             text                           not null,
    "ParagraphStyleId" bigint                         not null
        constraint "Styles_ParagraphStyle_Id_fk"
            references public."ParagraphStyle",
    "TextStyleId"      bigint                         not null
        constraint "Styles_TextStyle_Id_fk"
            references public."TextStyle"
);

alter table public."Styles"
    owner to postgres;

create unique index if not exists "Styles_Name_uindex"
    on public."Styles" ("Name");

create table if not exists public."DocumentsStyles"
(
    "Id"         bigint generated always as identity
        constraint "DocumentsStyles_pk"
            primary key,
    "DocumentId" uuid not null
        constraint "DocumentsStyles_Documents_Id_fk"
            references public."Documents",
    "StyleId"    uuid not null
        constraint "DocumentsStyles_Styles_Id_fk"
            references public."Styles"
);

alter table public."DocumentsStyles"
    owner to postgres;

