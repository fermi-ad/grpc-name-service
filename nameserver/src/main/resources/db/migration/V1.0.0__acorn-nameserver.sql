
    create sequence AlarmType_SEQ start with 1 increment by 50;

    create sequence channel_alarm_SEQ start with 1 increment by 50;

    create sequence Channel_SEQ start with 1 increment by 50;

    create sequence ChannelAccessControl_SEQ start with 1 increment by 50;

    create sequence ChannelTransform_SEQ start with 1 increment by 50;

    create sequence Device_SEQ start with 1 increment by 50;

    create sequence Location_SEQ start with 1 increment by 50;

    create sequence LocationType_SEQ start with 1 increment by 50;

    create sequence Node_SEQ start with 1 increment by 50;

    create sequence Role_SEQ start with 1 increment by 50;

    create table AlarmType (
        id bigint not null,
        description varchar(255),
        name varchar(255) unique,
        primary key (id)
    );

    create table Channel (
        device_id bigint not null,
        id bigint not null,
        description varchar(255),
        metadata varchar(255),
        name varchar(255) unique,
        primary key (id)
    );

    create table ChannelAlarm (
        alarm_type_id bigint not null,
        channel_id bigint not null,
        id bigint not null,
        trigger_condition varchar(255),
        primary key (id),
        unique (alarm_type_id, channel_id)
    );

    create table ChannelAccessControl (
        read boolean not null,
        write boolean not null,
        channel_id bigint not null,
        id bigint not null,
        role_id bigint not null,
        primary key (id)
    );

    create table ChannelTransform (
        channel_id bigint not null,
        id bigint not null,
        description varchar(255),
        name varchar(255) unique,
        transform varchar(255),
        primary key (id)
    );

    create table Device (
        id bigint not null,
        node_id bigint not null,
        description varchar(255),
        name varchar(255) unique,
        primary key (id)
    );

    create table Location (
        id bigint not null,
        location_type_id bigint,
        parent_location_id bigint,
        description varchar(255),
        name varchar(255) unique,
        primary key (id)
    );

    create table LocationType (
        id bigint not null,
        description varchar(255),
        name varchar(255) unique,
        primary key (id)
    );

    create table Node (
        id bigint not null,
        location_id bigint not null,
        description varchar(255),
        hostname varchar(255) unique,
        ip_address varchar(255),
        primary key (id)
    );

    create table Role (
        id bigint not null,
        description varchar(255),
        name varchar(255) unique,
        primary key (id)
    );

    alter table if exists Channel 
       add constraint FKsemie3xgis3q69rhigv8y2qln 
       foreign key (device_id) 
       references Device;

    alter table if exists channel_alarm 
       add constraint FKopcl5pa9y9h8tk9nufpx7ptxm 
       foreign key (alarm_type_id) 
       references AlarmType;

    alter table if exists channel_alarm 
       add constraint FKgvjgq3wttiulpuvg7s81kwkbs 
       foreign key (channel_id) 
       references Channel;

    alter table if exists ChannelAccessControl 
       add constraint FKh5owa27lx3fr7yxwq4f7l8rcx 
       foreign key (channel_id) 
       references Channel;

    alter table if exists ChannelAccessControl 
       add constraint FK9wh41r15mum5hv4i76obmvqga 
       foreign key (role_id) 
       references Role;

    alter table if exists ChannelTransform 
       add constraint FKf1kf0j4vukcstcmw8ik3soyse 
       foreign key (channel_id) 
       references Channel;

    alter table if exists Device 
       add constraint FKiw4ssv0383hx8u30p4h0clh7r 
       foreign key (node_id) 
       references Node;

    alter table if exists Location 
       add constraint FKpsrwdjrj682g76jj8hfvwt31d 
       foreign key (location_type_id) 
       references LocationType;

    alter table if exists Location 
       add constraint FK9ukkdd3ew9ydpdqcxs8nhlg6p 
       foreign key (parent_location_id) 
       references Location;

    alter table if exists Node 
       add constraint FKcj4av5k5htqi1vhdptpfacbw1 
       foreign key (location_id) 
       references Location;
       
