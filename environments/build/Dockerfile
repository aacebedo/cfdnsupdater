#This file is part of CFDNSUpdater.
#
#    CFDNSUpdater is free software: you can redistribute it and/or modify
#    it under the terms of the GNU General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.
#
#    CFDNSUpdater is distributed in the hope that it will be useful,
#    but WITHOUT ANY WARRANTY; without even the implied warranty of
#    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
#    GNU General Public License for more details.
#
#    You should have received a copy of the GNU General Public License
#    along with CFDNSUpdater.  If not, see <http://www.gnu.org/licenses/>.

from alpine
RUN apk update
RUN apk upgrade
RUN apk add go git python3 bash ruby-dev gcc make ruby ruby-bundler \
 g++ libstdc++ libffi-dev ruby-rdoc
RUN gem install fpm

ENV GOPATH /go
ENV GOBIN /go/bin/usr/local/bin
ENV CGO_ENABLED=0

ADD ./build.py /usr/bin
RUN chmod u+x /usr/bin/build.py

WORKDIR /out
VOLUME ["/out"]

#ENTRYPOINT ["build.py", "-p", "aacebedo/cfdnsupdater", "-n", "cfdnsupdater"]



