insert into alarmtype (id, name, description) 
values (nextval('alarmtype_SEQ'), 'Major', 'Major severity');
insert into alarmtype (id, name, description) 
values (nextval('alarmtype_SEQ'), 'Minor', 'Minor severity');

insert into locationtype (id, name, description) 
values (nextval('locationtype_SEQ'), 'Building', 'Building location');
insert into locationtype (id, name, description) 
values (nextval('locationtype_SEQ'), 'Area', 'Area location');
insert into locationtype (id, name, description) 
values (nextval('locationtype_SEQ'), 'Rack', 'Rack location');

insert into location (id, name, description, location_type_id) 
values (nextval('location_SEQ'), 'Building A', 'Building A description', 1);

insert into node (id, hostname, description, location_id, ip_address) 
values (nextval('node_SEQ'), 'Node A', 'Node A description', 1, '1.1.1.1');


insert into role(id, name, description) 
values (nextval('role_SEQ'), 'Admin', 'Admin role');
insert into role(id, name, description) 
values (nextval('role_SEQ'), 'Support', 'Support role');
insert into role(id, name, description) 
values (nextval('role_SEQ'), 'General', 'General role');

insert into device (id, name, description, node_id)
values(nextval('device_SEQ'), 'Device A', 'Device A description', 1);

insert into channel (id, name, description, device_id)
values(nextval('channel_SEQ'), 'Channel A', 'Channel A description', 1);