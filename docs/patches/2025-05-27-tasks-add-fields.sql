alter table public.companies
    add "createdAt" timestamp default now() not null;

alter table public.tasks
    add name text;

alter table public.tasks
    add url text;

alter table public.tasks
    add "createdAt" timestamp default now() not null;

alter table public.students
    add "createdAt" timestamp default now() not null;

alter table public.tasks
alter column deadline type text using deadline::text;