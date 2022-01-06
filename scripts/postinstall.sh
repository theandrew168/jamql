#!/bin/sh
set -e

# Create jamql group (if it doesn't exist)
if ! getent group jamql >/dev/null; then
    groupadd --system jamql
fi

# Create jamql user (if it doesn't exist)
if ! getent passwd jamql >/dev/null; then
    useradd                            \
        --system                       \
        --gid jamql                    \
        --shell /usr/sbin/nologin      \
        --comment "jamql feed reader"  \
        jamql
fi

# Update config file permissions (idempotent)
chown root:jamql /etc/jamql.conf
chmod 0640 /etc/jamql.conf

# Reload systemd to pickup jamql.service
systemctl daemon-reload
