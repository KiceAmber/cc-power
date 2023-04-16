create database `cc-power`;

use `cc-power`;

create table `user` (
    `id` bigint not null auto_increment primary key,
    `user_id` varchar(128) not null,
    `username` varchar(20) not null,
    `password` varchar(20) not null,
    `role` varchar(20) not null,
    `created_at` TimeStamp not null default current_timestamp
);

insert into `user`(`user_id`, `username`, `password`, `role`) values('10000', 'admin', '123456', '管理员');
insert into `user`(`user_id`, `username`, `password`, `role`) values('10001', 'tom', '123456', '购电用户');
