# Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
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

FROM azul/zulu-openjdk-alpine:17 AS jlink

RUN "$JAVA_HOME/bin/jlink" --compress=2 \
    --module-path /opt/java/openjdk/jmods \
    --add-modules java.base,java.compiler,java.datatransfer,jdk.crypto.ec,java.desktop,java.instrument,java.logging,java.management,java.naming,java.rmi,java.scripting,java.security.sasl,java.sql,java.transaction.xa,java.xml,jdk.unsupported \
    --output /jlinked

FROM python:3.10.4-alpine3.14

RUN apk update && apk add curl && apk add git

RUN apk add build-base

RUN pip install semgrep

# Create Folder for Rules semgrep
# START
RUN rm -rf /opt/rules-semgrep

RUN git clone https://github.com/jasayosep/rules-semgrep.git /opt/rules-semgrep
# END

COPY --from=jlink /jlinked /opt/jdk/

ENV JAVA_HOME=/opt/jdk

# Install Dependency Check OWASP
ENV VERSION_DEPENDENCY_CHECK=9.1.0
RUN curl -o /bin/dependency-check-release.zip -LO https://github.com/jeremylong/DependencyCheck/releases/download/v$VERSION_DEPENDENCY_CHECK/dependency-check-$VERSION_DEPENDENCY_CHECK-release.zip
RUN unzip /bin/dependency-check-release.zip -d /bin
RUN rm /bin/dependency-check-release.zip

# Download Trivy Latest
RUN curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /usr/local/bin latest

# Copy Database
RUN mkdir -p $HOME/.cache/trivy/db
RUN trivy image --download-db-only

#Copy Database JAVA
RUN mkdir -p $HOME/.cache/trivy/java-db
RUN trivy image --download-java-db-only
