CREATE TABLE public.customereventstore
(
  eventid UUID PRIMARY KEY NOT NULL,
  customerid UUID NOT NULL,
  eventdata JSONB DEFAULT '{}' NOT NULL,
  createdat TIMESTAMPTZ NOT NULL
);
CREATE UNIQUE INDEX customereventstore_eventid_uindex ON public.customereventstore (eventid);