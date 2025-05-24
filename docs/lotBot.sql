-- =============================================================================
-- Diagram Name: apisrv
-- Created on: 24.05.2025 4:39:35
-- Diagram Version:
-- =============================================================================

CREATE TABLE "companies" (
                             "companyId" SERIAL NOT NULL,
                             "name" text NOT NULL,
                             "tgId" int8 NOT NULL,
                             "inn" text NOT NULL,
                             "scope" text NOT NULL,
                             "userName" varchar(128) NOT NULL,
                             "phone" text NOT NULL,
                             "statusId" int4 NOT NULL,
                             PRIMARY KEY("companyId")
);



CREATE TABLE "tasks" (
                         "taskId" SERIAL NOT NULL,
                         "companyId" int4 NOT NULL,
                         "name" text,
                         "scope" text NOT NULL,
                         "description" text NOT NULL,
                         "link" text NOT NULL,
                         "deadline" text NOT NULL,
                         "contactSlot" text NOT NULL,
                         "statusId" int4 NOT NULL,
                         "studentId" int4,
                         "budget" float8 NOT NULL,
                         "yougileId" text,
                         "url" text,
                         "createdAt" timestamp NOT NULL DEFAULT now(),
                         PRIMARY KEY("taskId")
);

CREATE TABLE "students" (
                            "studentId" SERIAL NOT NULL,
                            "tgId" int8 NOT NULL,
                            "name" text NOT NULL,
                            "birthday" text NOT NULL,
                            "city" text NOT NULL,
                            "scope" text NOT NULL,
                            "email" text NOT NULL,
                            "statusId" int4 NOT NULL,
                            PRIMARY KEY("studentId")
);


ALTER TABLE "tasks" ADD CONSTRAINT "Ref_tasks_to_companies" FOREIGN KEY ("companyId")
    REFERENCES "companies"("companyId")
        MATCH SIMPLE
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
    NOT DEFERRABLE;

ALTER TABLE "tasks" ADD CONSTRAINT "Ref_tasks_to_students" FOREIGN KEY ("studentId")
    REFERENCES "students"("studentId")
        MATCH SIMPLE
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
    NOT DEFERRABLE;


