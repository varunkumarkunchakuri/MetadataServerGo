# MetadataServerGo

## 1 Start the server  by executing the command

"go run server.go"

## 2 Run the tests by executing the command

"go test"

## 3 Dashboards required for operationalizing the service
<ol>
<li>Request rate and latency over time</li>
<li>Error rate over time</li>
<li>Number of records stored over time</li>
<li>Distribution of request latencies</li>
<li>Distribution of request sizes</li>
<li>Distribution of error types</li>
<li>Distribution of response codes</li>
</ol>

## 4 Alerts required for operationalizing the service
<ol>
<li>High error rate</li>
<li>High latency</li>
<li>Low disk space or memory based on the implementation</li>
<li>High number of simultaneous requests(ddos)</li>
<li>High number of records stored</li>
</ol>

