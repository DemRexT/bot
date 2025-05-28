alter table tasks
    add name text;

alter table tasks
alter column deadline type text using deadline::text;

alter table tasks
    add url text;

alter table tasks
    add "createdAt" timestamp default now() not null;