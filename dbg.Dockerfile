# Copyright The KubeVault Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM {ARG_FROM}

ENV DEBIAN_FRONTEND noninteractive
ENV DEBCONF_NONINTERACTIVE_SEEN true

RUN set -x \
  && apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates tzdata locales e2fsprogs mount \
  && rm -rf /var/lib/apt/lists/* /usr/share/doc /usr/share/man /tmp/* \
  && localedef -i en_US -c -f UTF-8 -A /usr/share/locale/locale.alias en_US.UTF-8 \
  && echo 'Etc/UTC' > /etc/timezone && dpkg-reconfigure tzdata

ENV TZ     :/etc/localtime
ENV LANG   en_US.utf8
ENV LC_ALL en_US.UTF-8

ADD bin/{ARG_OS}_{ARG_ARCH}/{ARG_BIN} /{ARG_BIN}

ENTRYPOINT ["/{ARG_BIN}"]
