# Các Metrics Mặc Định từ grpc_prometheus

| Tên Metrics                          | Loại     | Mô tả                                                                 |
|-------------------------------------|----------|----------------------------------------------------------------------|
| grpc_server_started_total           | Counter  | Số lượng RPC được bắt đầu (theo method và service)                  |
| grpc_server_handled_total           | Counter  | Số lượng RPC đã xử lý xong, phân theo mã lỗi (OK, CANCELLED,...)   |
| grpc_server_msg_received_total      | Counter  | Tổng số messages được nhận từ client (áp dụng cho stream RPC)      |
| grpc_server_msg_sent_total          | Counter  | Tổng số messages được gửi đến client (áp dụng cho stream RPC)      |
| grpc_server_handling_seconds_count  | Histogram| Số lượng RPC được đo thời gian xử lý                                |
| grpc_server_handling_seconds_sum    | Histogram| Tổng thời gian xử lý RPC (đơn vị giây)                              |
| grpc_server_handling_seconds_bucket | Histogram| Phân phối thời gian xử lý theo buckets (dùng để tính percentile)   |


# Sau khi gọi EnableHandlingTimeHistogram()

```go
grpc_prometheus.EnableHandlingTimeHistogram()
```
Sẽ bật histogram cho thời gian xử lý request, mặc định là không bật. Khi bật rồi, các metrics sau được export thêm:

| Tên Metrics	| Ý nghĩa
|-------------------------------------|----------|
| grpc_server_handling_seconds_bucket	| Biểu đồ thời gian xử lý các RPC chia theo buckets (ví dụ: <0.01s, <0.05s, <0.1s,...) |
| grpc_server_handling_seconds_sum	    | Tổng thời gian xử lý RPC |
| grpc_server_handling_seconds_count	| Số lượng RPC được đo thời gian |

# Các Label kèm theo
Mỗi metric đều có các label (dimension) để bạn có thể lọc hoặc phân nhóm trên Grafana:

| Label	| Mô tả
| --- | --- |
| grpc_type	    | Loại RPC (unary, client_stream, server_stream, bidi_stream) |
| grpc_service	| Tên service gRPC (vd: helloworld.Greeter) |
| grpc_method	| Tên method RPC (vd: SayHello) |
| grpc_code	    | Trạng thái kết quả (vd: OK, CANCELLED, INTERNAL,...) |



Trên Grafana bạn có thể tạo các biểu đồ như:
- Request Rate (QPS):
```bash
rate(grpc_server_started_total[1m])
```

- Error Rate:
```bash
rate(grpc_server_handled_total{grpc_code!="OK"}[1m])
```

- Latency (95th percentile):
```bash
histogram_quantile(0.95, sum(rate(grpc_server_handling_seconds_bucket[5m])) by (le))
```

- Số lượng request theo method:
```bash
rate(grpc_server_started_total[1m])
  by (grpc_service, grpc_method)
```

- Top RPCs theo thời gian xử lý trung bình:
```bash
sum(rate(grpc_server_handling_seconds_sum[1m])) by (grpc_service, grpc_method)
/ 
sum(rate(grpc_server_handling_seconds_count[1m])) by (grpc_service, grpc_method)
```


