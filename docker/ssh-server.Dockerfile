FROM dunedaq/sl7-minimal

# an up to date python setup and ssh server
RUN yum install -y openssh-server python3-pip && \
    python3 -m pip install --upgrade pip && \
    yum clean all

# setup pre-defined keys for convenience
ADD ssh.key.pub /root/.ssh/authorized_keys
ADD ssh.key /root/.ssh/id_rsa
ADD sshd_config /etc/ssh/sshd_config
RUN chmod -R 700 /root/.ssh && \
    chmod 400 /root/.ssh/authorized_keys && \
    chmod 600 /root/.ssh/id_rsa

RUN /usr/sbin/sshd-keygen || true

ENTRYPOINT [ "/usr/sbin/sshd" ]
CMD [ "-D" ]