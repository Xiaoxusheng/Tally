create table collect_basic
(
    id            bigint unsigned auto_increment
        primary key,
    created_at    datetime(3) null,
    updated_at    datetime(3) null,
    deleted_at    datetime(3) null,
    identity      varchar(36) not null comment '''唯一标识''',
    user_identity varchar(36) not null comment '''用户唯一标识''',
    collect_id    varchar(36) not null comment '''被关注用户唯一标识''',
    blog_id       varchar(36) not null comment '''blog唯一标识''',
    constraint blog_id
        unique (blog_id),
    constraint collect_id
        unique (collect_id),
    constraint identity
        unique (identity),
    constraint user_identity
        unique (user_identity)
);

create index idx_collect_basic_deleted_at
    on collect_basic (deleted_at);

create table comment_basic
(
    id            bigint unsigned auto_increment
        primary key,
    created_at    datetime(3) null,
    updated_at    datetime(3) null,
    deleted_at    datetime(3) null,
    identity      varchar(36) not null comment '''唯一标识''',
    user_identity varchar(36) not null comment '''用户唯一标识''',
    blog_id       varchar(36) not null comment '''blog唯一标识''',
    parent_id     varchar(36) not null comment '''父评论ID''',
    text          varchar(36) not null comment '''评论内容''',
    ip           varchar(36) not null comment '''ip地址''',
    is_top       tinyint(1) null comment '''是否为父评论''',
    violate_rule tinyint(1) null comment '''评论内容是否违规''',
    constraint identity
        unique (identity),
    constraint user_identity
        unique (user_identity)
);

create index idx_comment_basic_deleted_at
    on comment_basic (deleted_at);

create table follow_basic
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(3) null,
    updated_at datetime(3) null,
    deleted_at datetime(3) null,
    identity   varchar(36) not null comment '''关注记录唯一标识''',
    user_id    varchar(36) not null comment '''用户唯一标识''',
    follow_id  varchar(36) not null comment '''关注用户唯一标识''',
    constraint identity
        unique (identity)
);

create index idx_follow_basic_deleted_at
    on follow_basic (deleted_at);

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
    collect tinyint(1) null comment '''是否收藏''',
    remark        varchar(255) not null comment '''备注''',
    category      int          not null comment '''类别''',
    constraint identity
        unique (identity)
);

create index idx_tally_basic_deleted_at
    on tally_basic (deleted_at);

create table tally_blog_basic
(
    id            bigint unsigned auto_increment
        primary key,
    created_at    datetime(3)   null,
    updated_at    datetime(3)   null,
    deleted_at    datetime(3)   null,
    identity      varchar(36)   not null comment '''唯一标识''',
    user_identity varchar(36)   not null comment '''用户唯一标识''',
    img_url       varchar(1000) null comment '''图片url,为空表示没有图片''',
    text          varchar(2000) null comment '''文本内容'' ',
    is_hide      int null comment '''文章是否私密''',
    likes         int           null comment '''点赞数量''',
    ip           varchar(64) not null comment '''IP地址''',
    violate_rule tinyint(1) not null comment '''评论内容是否违规''',
    constraint identity
        unique (identity)
);

create index idx_tally_blog_basic_deleted_at
    on tally_blog_basic (deleted_at);

create table user_basic
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(3) null,
    updated_at datetime(3) null,
    deleted_at datetime(3) null,
    username  varchar(10) null comment '''用户名''',
    account   int         not null comment '''账号''',
    password  varchar(36) not null comment '''密码''',
    phone     varchar(11) not null comment '''手机号''',
    identity  varchar(36) not null comment '''唯一标识''',
    github_id varchar(36) not null comment '''Github账号''',
    status    int null comment '''0表示正常, 1表示封禁''',
    is_hide   int null comment '''是否隐私账号''',
    ip        varchar(64) not null comment '''IP地址''',
    constraint account
        unique (account),
    constraint identity
        unique (identity),
    constraint phone
        unique (phone)
);

create index idx_user_basic_deleted_at
    on user_basic (deleted_at);

