/*
 * An implementation of key value pair (KVP) functionality for Linux.
 *
 *
 * Copyright (C) 2010, Novell, Inc.
 * Author : K. Y. Srinivasan <ksrinivasan@novell.com>
 *
 * This program is free software; you can redistribute it and/or modify it
 * under the terms of the GNU General Public License version 2 as published
 * by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, GOOD TITLE or
 * NON INFRINGEMENT.  See the GNU General Public License for more
 * details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA 02110-1301 USA.
 *
 */

#include <sys/poll.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>
#include <fcntl.h>
#include <errno.h>
#include <arpa/inet.h>
#include <linux/hyperv.h>

#define TIMEOUT 5000

int do_read_kvp_data(FILE *dest)
{
	int kvp_fd = -1, len;
	int error;
	struct pollfd pfd;
	char    *p;
	struct hv_kvp_msg hv_msg[1];
	char	*key_value;
	char	*key_name;
	int	op;
	int	pool;
	char	*if_name;
	struct hv_kvp_ipaddr_value *kvp_ip_val;
	int in_hand_shake;
	uint32_t zero = 0;
	size_t l;

	if (kvp_fd != -1)
		close(kvp_fd);
	in_hand_shake = 1;

	kvp_fd = open("/dev/vmbus/hv_kvp", O_RDWR | O_CLOEXEC);
	if (kvp_fd < 0)
		return kvp_fd;

	hv_msg->kvp_hdr.operation = KVP_OP_REGISTER1;
	len = write(kvp_fd, hv_msg, sizeof(struct hv_kvp_msg));
	if (len != sizeof(struct hv_kvp_msg)) {
		close(kvp_fd);
		return -1;
	}

	pfd.fd = kvp_fd;

	while (1) {
		int how_many;

		pfd.events = POLLIN;
		pfd.revents = 0;

                how_many = poll(&pfd, 1, TIMEOUT);
		if (how_many < 0) {
			if (errno == EINVAL) {
				close(kvp_fd);
				return -1;
			}
			else
				continue;
		}
        	if (how_many == 0)
			return 0;

		len = read(kvp_fd, hv_msg, sizeof(struct hv_kvp_msg));
		if (len != sizeof(struct hv_kvp_msg)) {
			return -1;
		}

		op = hv_msg->kvp_hdr.operation;
		pool = hv_msg->kvp_hdr.pool;
		hv_msg->error = HV_S_OK;

		if ((in_hand_shake) && (op == KVP_OP_REGISTER1)) {
			in_hand_shake = 0;
			continue;
		}

		if (op == KVP_OP_SET) {
			l = fwrite(&hv_msg->body.kvp_set.data.key_size, 1, sizeof(hv_msg->body.kvp_set.data.key_size), dest);
			if (l != sizeof(hv_msg->body.kvp_set.data.key_size)) {
				close(kvp_fd);
				return -1;
			}

			l = fwrite(&hv_msg->body.kvp_set.data.key, 1, hv_msg->body.kvp_set.data.key_size, dest);
			if (l != hv_msg->body.kvp_set.data.key_size) {
				close(kvp_fd);
				return -1;
			}

			l = fwrite(&hv_msg->body.kvp_set.data.value_size, 1, sizeof(hv_msg->body.kvp_set.data.value_size), dest);
			if (l != sizeof(hv_msg->body.kvp_set.data.value_size)) {
				close(kvp_fd);
				return -1;
			}

			l = fwrite(&hv_msg->body.kvp_set.data.value, 1, hv_msg->body.kvp_set.data.value_size, dest);
			if (l != hv_msg->body.kvp_set.data.value_size) {
				close(kvp_fd);
				return -1;
			}
		}

		len = write(kvp_fd, hv_msg, sizeof(struct hv_kvp_msg));
		if (len != sizeof(struct hv_kvp_msg)) {
			close(kvp_fd);
			return -1;
		}
	}
	close(kvp_fd);

	l = fwrite(&zero, 1, sizeof(zero), dest);
	if (l != sizeof(zero)) {
		return -1;
	}
	l = fwrite(&zero, 1, sizeof(zero), dest);
	if (l != sizeof(zero)) {
		return -1;
	}

	return 0;
}

char *read_kvp_data(void)
{
	char *buf = NULL;
	size_t len = 0;
	FILE *f;

	f = open_memstream(&buf, &len);

	if (do_read_kvp_data(f) < 0)
		return NULL;
	fclose(f);
        
	return buf;
}
