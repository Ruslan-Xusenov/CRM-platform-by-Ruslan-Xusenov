# ==============================================================
# Asterisk 20 PBX — Dockerfile
# ==============================================================
FROM debian:bookworm-slim

ENV DEBIAN_FRONTEND=noninteractive

# Install Asterisk 20 and required modules
RUN apt-get update && apt-get install -y --no-install-recommends \
    asterisk \
    asterisk-modules \
    asterisk-core-sounds-en \
    asterisk-moh-opsound-wav \
    curl \
    jq \
    ca-certificates \
    python3 \
    python3-pip \
    && rm -rf /var/lib/apt/lists/*

# Install MinIO client for recording uploads
RUN curl -fsSL https://dl.min.io/client/mc/release/linux-amd64/mc -o /usr/local/bin/mc && \
    chmod +x /usr/local/bin/mc

# Copy custom configuration
COPY conf/ /etc/asterisk/

# Copy recording upload script
COPY scripts/upload_recording.sh /usr/local/bin/upload_recording.sh
RUN chmod +x /usr/local/bin/upload_recording.sh

# Create recording directory
RUN mkdir -p /var/spool/asterisk/recording && \
    chown -R asterisk:asterisk /var/spool/asterisk/recording

EXPOSE 5060/udp 5060/tcp 8088 8089 10000-10100/udp

# Run as asterisk user
USER asterisk

CMD ["asterisk", "-fvvv"]
