CREATE TABLE document (
       id VARCHAR(50) PRIMARY KEY,
       read_only_id VARCHAR(50) UNIQUE NOT NULL,
       is_finished BOOLEAN DEFAULT FALSE,
       created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE document_operation (
       id VARCHAR(50) NOT NULL,
       document_id VARCHAR(50) NOT NULL,
       editor_id VARCHAR(50) NOT NULL,
       position_start BIGINT NOT NULL,
       position_end BIGINT NOT NULL,
       val VARCHAR(10000),
       op VARCHAR(20) NOT NULL,
       operation_time TIMESTAMP,

       PRIMARY KEY (id, document_id)
);
