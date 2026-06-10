from google.cloud import bigquery

client = bigquery.Client.from_service_account_json('arquivo.json')

query = """
SELECT *
FROM `conciliation.gold_flytura`

where (emission_date>='2026-04-01' and emission_date<='2026-06-09') 
LIMIT 10
"""

query_job = client.query(query)
results = query_job.result()

for row in results:
    print(dict(row))