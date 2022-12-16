#!/sbin/openrc-run

name=$RC_SVCNAME
supervisor="supervise-daemon"
command="/usr/local/bin/agent"
pidfile="/run/agent.pid"
command_user="root:root"

depend() {
	after net
}
