# Recgo-Engine

`Recgo-Engine` is a high-performance recommendation engine written in Go.  
It exposes HTTP APIs for retrieving **feeds** and **related item** recommendations, based on configurable pipelines.

---

## API Endpoints

| Method | Path       | Description                                       |
|--------|-----------|---------------------------------------------------|
| POST   | `/feeds`  | Get feed stream recommendations for a user         |
| POST   | `/related`| Get related item recommendations (detail page etc) |

---

## Request Structure

```json
{
  "trace_id": "abcd-1234-efgh-5678",
  "user_id": "u12345",
  "pipeline": "main_feed",
  "relate_id": "item987",
  "count": 10,
  "context": { "device": {...}, "version": {...} },
  "features": { "vip_level": {...} }
}
```

**Fields**

| Name       | Type   | Required | Description |
|------------|--------|----------|-------------|
| `trace_id` | string | optional | Unique request ID for tracking/debugging |
| `user_id`  | string | **yes**  | Unique user ID |
| `pipeline` | string | **yes**  | Pipeline name to execute; must match configured pipeline in engine |
| `relate_id`| string | optional | Target item ID for `related` recommendations |
| `count`    | int32  | optional | Number of items to return; default pipeline config applies if omitted |
| `context`  | object | optional | Contextual features (e.g., device, region) |
| `features` | object | optional | External features provided by caller |

---

## Response Structure

```json
{
  "code": 0,
  "message": "success",
  "trace_id": "abcd-1234-efgh-5678",
  "user_id": "u12345",
  "pipeline": "main_feed",
  "items": [
    {
      "item": "item001",
      "channels": ["recallA", "recallB"],
      "reasons": ["similar to watched item", "popular in category"]
    },
    {
      "item": "item002",
      "channels": ["recallC"],
      "reasons": ["user preference match"]
    }
  ],
  "count": 2
}
```

**Fields**

| Name        | Type     | Description |
|-------------|----------|-------------|
| `code`      | int      | Business status code (0 = success, non-zero = error) |
| `message`   | string   | Status or error description |
| `trace_id`  | string   | Echo from request for tracking |
| `user_id`   | string   | Echo from request |
| `pipeline`  | string   | Pipeline actually executed |
| `items`     | array    | List of recommended items |
| `count`     | int      | Number of items returned |

**ItemInfo fields**:
- `item`: Item ID  
- `channels`: Names of recall channels contributing to recommendation  
- `reasons`: Explanation texts for why item is recommended  

---

## Error Codes

| Code  | Meaning                                  |
|-------|------------------------------------------|
| 0     | Success                                  |
| -1    | Pipeline not found (`pipeline` unknown)  |

---

## Example Requests

### `/feeds`
```bash
curl -X POST "http://localhost:8080/feeds" \
     -H "Content-Type: application/json" \
     -d '{
           "trace_id": "req-0001",
           "user_id": "u001",
           "pipeline": "main_feed",
           "count": 5
         }'
```

### `/related`
```bash
curl -X POST "http://localhost:8080/related" \
     -H "Content-Type: application/json" \
     -d '{
           "trace_id": "req-0002",
           "user_id": "u002",
           "pipeline": "related_items",
           "relate_id": "item123",
           "count": 5
         }'
```

---

## Call Flow

```mermaid
flowchart TD
    A[POST /feeds or /related] --> B[Validate Request Fields]
    B --> C[Lookup Pipeline by name]
    C -->|Found| D[runPipeline(userCtx, pipeline)]
    C -->|Not Found| E[Return error: code=-1]
    D --> F[Pipeline Executes: Recaller→Filter→Ranker]
    F --> G[Collect Items & Build Response]
    G --> H[Return JSON to Client]
```

---

## Notes
- `pipeline` is **mandatory** for both `/feeds` and `/related` requests.
- `/related` requests may include `relate_id` for context-specific recommendations.
- `trace_id` helps track logs and metrics across services.

---

## License
Part of the **recgo-engine** project. See LICENSE for details.
