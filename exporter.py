from flask import Flask, Response
import openstack
import os
import configparser

app = Flask(__name__)

def get_openstack_connection(config):
    conn = openstack.connect(
        auth_url=config.get("Global", "auth-url"),
        project_id=config.get("Global", "tenant-id"),
        username=config.get("Global", "username").strip('"'),
        password=config.get("Global", "password").strip('"'),
        user_domain_id=config.get("Global", "domain-id"),
        project_domain_id=config.get("Global", "domain-id"),
        region_name=config.get("Global", "region"),
        verify=config.getboolean("Global", "tls-insecure") is False
    )
    return conn

def collect_metrics(conn):
    # Collect compute limits
    compute_limits = conn.get_compute_limits().to_dict()
    compute_metrics = [
        f'openstack_compute_{key} {value}'
        for key, value in compute_limits.items()
        if isinstance(value, (int, float))
    ]

    # Collect volume limits
    volume_limits = conn.get_volume_limits().to_dict()
    volume_metrics = [
        f'openstack_volume_{key} {value}'
        for key, value in volume_limits['absolute'].items()
        if isinstance(value, (int, float))
    ]

    # Combine all metrics
    all_metrics = compute_metrics + volume_metrics
    return "\n".join(all_metrics)

@app.route('/metrics')
def metrics():
    config = configparser.ConfigParser()
    config.read('cloud.conf')  # Replace with the actual path to your config file
    conn = get_openstack_connection(config)
    metrics_data = collect_metrics(conn)
    return Response(metrics_data, mimetype='text/plain')

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=9183)
