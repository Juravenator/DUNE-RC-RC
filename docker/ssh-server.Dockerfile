FROM dunedaq/sl7-minimal

RUN yum install -y openssh-server python3-pip vim && \
    python3 -m pip install --upgrade pip && \
    yum clean all

ADD ssh.key.pub /root/.ssh/authorized_keys
ADD ssh.key /root/.ssh/id_rsa
ADD sshd_config /etc/ssh/sshd_config
RUN chmod -R 700 /root/.ssh && \
    chmod 400 /root/.ssh/authorized_keys && \
    chmod 600 /root/.ssh/id_rsa

RUN /usr/sbin/sshd-keygen || true

ENTRYPOINT [ "/usr/sbin/sshd" ]
CMD [ "-D" ]