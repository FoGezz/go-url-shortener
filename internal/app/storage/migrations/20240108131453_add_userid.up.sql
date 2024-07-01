alter table links
    add user_uuid uuid null;

create index links_user_uuid_index
    on links (user_uuid);