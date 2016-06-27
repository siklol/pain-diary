CREATE TABLE public.patienteventstore
(
  eventid UUID PRIMARY KEY NOT NULL,
  patientid UUID NOT NULL,
  eventdata JSONB DEFAULT '{}' NOT NULL,
  createdat TIMESTAMPTZ NOT NULL
);
CREATE UNIQUE INDEX patienteventstore_eventid_uindex ON public.patienteventstore (eventid);