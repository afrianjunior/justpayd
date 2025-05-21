```dbml
Table workers {
  id integer [primary key, autoincrement]
  name text [not null]
  email text [unique, not null]
  created_at timestamp [default: `CURRENT_TIMESTAMP`]
}

Table shifts {
  id integer [primary key, autoincrement]
  date date [not null]
  start_time time [not null]
  end_time time [not null]
  role text [not null]
  location text
  created_at timestamp [default: `CURRENT_TIMESTAMP`]
}

Table shift_requests {
  id integer [primary key, autoincrement]
  worker_id integer [not null, ref: > workers.id]
  shift_id integer [not null, ref: > shifts.id]
  status text [not null]
  requested_at timestamp [default: `CURRENT_TIMESTAMP`]
  
  indexes {
    (worker_id, shift_id) [unique]
  }
}

Table assignments {
  id integer [primary key, autoincrement]
  shift_id integer [not null, unique, ref: > shifts.id]
  worker_id integer [not null, ref: > workers.id]
  assigned_at timestamp [default: `CURRENT_TIMESTAMP`]
}
```

You can use this DBML format with tools like:
- https://dbdiagram.io/
- https://dbdocs.io/
- https://holistics.io/dbml/

Just copy and paste the content above into these tools to visualize the database schema. 