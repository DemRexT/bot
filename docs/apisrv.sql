-- =============================================================================
-- Diagram Name: apisrv
-- Created on: 06.04.2025 18:02:07
-- Diagram Version: 
-- =============================================================================

CREATE TABLE "companies" (
	"companyId" SERIAL NOT NULL,
	"name" text NOT NULL,
	"tgId" text NOT NULL,
	"inn" int4 NOT NULL,
	"scope" text NOT NULL,
	"userName" name NOT NULL,
	"phone" int4 NOT NULL,
	"statusId" int4 NOT NULL,
	PRIMARY KEY("companyId")
);

CREATE TABLE "tasks" (
	"taskId" SERIAL NOT NULL,
	"companyId" int4 NOT NULL,
	"scope" text NOT NULL,
	"description" text NOT NULL,
	"images" varchar(128)[] NOT NULL,
	"link" text NOT NULL,
	"deadline" date NOT NULL,
	"contactSlot text" text NOT NULL,
	"statusId" int4 NOT NULL,
	"studentId" int4,
	PRIMARY KEY("taskId")
);

CREATE TABLE "students" (
	"studentId" SERIAL NOT NULL,
	"tgId" text NOT NULL,
	"name" text NOT NULL,
	"birthday" date NOT NULL,
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


