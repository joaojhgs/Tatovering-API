
CREATE TABLE agendamentos
(
  id                int      NOT NULL GENERATED ALWAYS AS IDENTITY UNIQUE,
  servico_id        int      NOT NULL,
  cliente_id        INT      NOT NULL,
  tatuador_id       int      NOT NULL,
  duracao           int      NOT NULL,
  valor             float    NOT NULL,
  status            enum    ,
  observacao        string   NOT NULL,
  data_inicio       datetime,
  data_criacao      datetime NOT NULL DEFAULT now(),
  data_atualizacao  datetime NOT NULL DEFAULT now(),
  data_cancelamento datetime,
  PRIMARY KEY (id)
);

COMMENT ON COLUMN agendamentos.servico_id IS 'ex: sessao, flash day, promocao,';

COMMENT ON COLUMN agendamentos.cliente_id IS 'id cliente';

COMMENT ON COLUMN agendamentos.duracao IS 'duracao / horas';

COMMENT ON COLUMN agendamentos.status IS 'ex: agendado, confirmado, cancelado, concluido';

COMMENT ON COLUMN agendamentos.data_inicio IS 'data/hora agendamento';

CREATE TABLE estudio
(
  id                    int      NOT NULL GENERATED ALWAYS AS IDENTITY UNIQUE,
  proprietario_id       INT      NOT NULL,
  nome                  text     NOT NULL,
  email                 varchar  DEFAULT NULL,
  horario_funcionamento jsonb    NOT NULL,
  endereco              text     DEFAULT NULL,
  localizacao           point    DEFAULT NULL,
  telefone              varchar  DEFAULT NULL,
  descricao             text    ,
  taxa_agendamento      money    DEFAULT NULL,
  data_criacao          datetime NOT NULL DEFAULT NOW(),
  data_atualizacao      datetime NOT NULL DEFAULT NOW(),
  data_exclusao         datetime DEFAULT NULL,
  dias_funcionamento    jsonb    NOT NULL,
  PRIMARY KEY (id)
);

COMMENT ON COLUMN estudio.email IS 'email profissional';

COMMENT ON COLUMN estudio.horario_funcionamento IS 'horarios de funcionamento do estudio';

COMMENT ON COLUMN estudio.endereco IS 'endereco';

COMMENT ON COLUMN estudio.localizacao IS 'localizacao em lat long';

COMMENT ON COLUMN estudio.taxa_agendamento IS 'O valor minimo que o tatuador cobra para agendar pela plataforma';

CREATE TABLE servico
(
  id               int      NOT NULL UNIQUE,
  tipo             text     NOT NULL,
  cliente_id       INT      NOT NULL,
  tatuador_id      int      NOT NULL,
  tatuagem_id      int      NOT NULL,
  descricao        string   NOT NULL,
  duracao          int      NOT NULL,
  preco            int     ,
  data_criacao     datetime DEFAULT NOW(),
  data_atualizacao datetime DEFAULT NOW(),
  data_exclusao    datetime,
  PRIMARY KEY (id)
);

CREATE TABLE tatuadores
(
  id               int       NOT NULL GENERATED ALWAYS AS IDENTITY,
  estudio_id       int       NOT NULL,
  experiencia      int       NOT NULL DEFAULT 1,
  estilo_tatuagem  text     ,
  status           text     ,
  tipo             text     ,
  redes_sociais    jsonb    ,
  estudio_id                ,
  data_criacao     timestamp DEFAULT NOW(),
  data_atualizacao timestamp DEFAULT NOW(),
  PRIMARY KEY (id)
);

COMMENT ON COLUMN tatuadores.experiencia IS 'meses';

CREATE TABLE tatuagens
(
  id               int    NOT NULL GENERATED ALWAYS AS IDENTITY UNIQUE,
  cliente_id       INT   ,
  tatuador_id      int    NOT NULL UNIQUE,
  agendamento_id   int    UNIQUE,
  preco            float ,
  desenho          string UNIQUE,
  tamanho          int   ,
  cor              enum  ,
  estilo           enum  ,
  data_criacao     date   DEFAULT NOW(),
  data_atualizacao date   DEFAULT NOW(),
  data_exclusao    date   DEFAULT NULL,
  PRIMARY KEY (id)
);

COMMENT ON COLUMN tatuagens.tamanho IS 'tamanho em cm';

CREATE TABLE usuarios
(
  id               INT       GENERATED ALWAYS AS IDENTITY,
  nome             text      NOT NULL,
  email            text     ,
  telefone_celular text     ,
  cpf              text     ,
  rg               text     ,
  data_nascimento  timestamp,
  status           text     ,
  experiencia      integer  ,
  endereco         text     ,
  data_criacao     timestamp DEFAULT NOW(),
  data_atualizacao timestamp DEFAULT NOW(),
  data_exclusao    timestamp,
  PRIMARY KEY (id)
);

ALTER TABLE tatuagens
  ADD CONSTRAINT FK_agendamentos_TO_tatuagens
    FOREIGN KEY (agendamento_id)
    REFERENCES agendamentos (id);

ALTER TABLE servico
  ADD CONSTRAINT FK_tatuagens_TO_servico
    FOREIGN KEY (tatuagem_id)
    REFERENCES tatuagens (id);

ALTER TABLE tatuagens
  ADD CONSTRAINT FK_usuarios_TO_tatuagens
    FOREIGN KEY (cliente_id)
    REFERENCES usuarios (id);

ALTER TABLE agendamentos
  ADD CONSTRAINT FK_usuarios_TO_agendamentos
    FOREIGN KEY (cliente_id)
    REFERENCES usuarios (id);

ALTER TABLE servico
  ADD CONSTRAINT FK_usuarios_TO_servico
    FOREIGN KEY (cliente_id)
    REFERENCES usuarios (id);

ALTER TABLE estudio
  ADD CONSTRAINT FK_usuarios_TO_estudio
    FOREIGN KEY (proprietario_id)
    REFERENCES usuarios (id);

ALTER TABLE tatuadores
  ADD CONSTRAINT FK_estudio_TO_tatuadores
    FOREIGN KEY (estudio_id)
    REFERENCES estudio (id);

ALTER TABLE tatuagens
  ADD CONSTRAINT FK_tatuadores_TO_tatuagens
    FOREIGN KEY (tatuador_id)
    REFERENCES tatuadores (id);

ALTER TABLE agendamentos
  ADD CONSTRAINT FK_tatuadores_TO_agendamentos
    FOREIGN KEY (tatuador_id)
    REFERENCES tatuadores (id);

ALTER TABLE servico
  ADD CONSTRAINT FK_tatuadores_TO_servico
    FOREIGN KEY (tatuador_id)
    REFERENCES tatuadores (id);

ALTER TABLE agendamentos
  ADD CONSTRAINT FK_servico_TO_agendamentos
    FOREIGN KEY (servico_id)
    REFERENCES servico (id);
