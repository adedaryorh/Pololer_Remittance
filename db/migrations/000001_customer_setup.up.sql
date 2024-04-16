CREATE TABLE "customer" (
                         "id" BIGSERIAL PRIMARY KEY,
                         "email" varchar(256) UNIQUE NOT NULL,
                         "hashed_password" varchar(256) NOT NULL,
                         "username" varchar(256) UNIQUE NOT NULL,
                         "firstname" varchar(256) UNIQUE NOT NULL,
                         "lastname" varchar(256) UNIQUE NOT NULL,
                         "gender" varchar(256) UNIQUE NOT NULL,
                         "state_of_origin" varchar(256) UNIQUE NOT NULL,
                         "created_at" timestamptz NOT NULL DEFAULT now(),
                         "updated_at" timestamptz NOT NULL DEFAULT now()
);





