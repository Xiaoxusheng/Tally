create table kind_basic
(
    id            bigint unsigned auto_increment
        primary key,
    created_at    datetime(3)  null,
    updated_at    datetime(3)  null,
    deleted_at    datetime(3)  null,
    name          varchar(255) not null,
    serial_number bigint       not null,
    constraint name
        unique (name),
    constraint serial_number
        unique (serial_number)
);

create index idx_kind_basic_deleted_at
    on kind_basic (deleted_at);

create table tally_basic
(
    id            bigint unsigned auto_increment
        primary key,
    created_at    datetime(3)  null,
    updated_at    datetime(3)  null,
    deleted_at    datetime(3)  null,
    identity      varchar(36)  not null comment '''唯一标识''',
    user_identity varchar(36)  not null comment '''用户唯一标识''',
    kind          bigint       null comment '''收入支出种类''',
    money         double       null comment '''金额''',
    remark        varchar(255) not null comment '''备注''',
    category      int          not null comment '''类别''',
    collect       tinyint(1)   null comment '''是否收藏''',
    constraint identity
        unique (identity)
);

create table user_basic
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(3) null,
    updated_at datetime(3) null,
    deleted_at datetime(3) null,
    username   varchar(10) not null,
    password   varchar(36) not null,
    phone      varchar(11) not null,
    identity   varchar(36) not null,
    ip         varchar(64) not null,
    constraint identity
        unique (identity),
    constraint phone
        unique (phone),

);

create index idx_user_basic_deleted_at
    on user_basic (deleted_at);

