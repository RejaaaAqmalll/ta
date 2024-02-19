QUERY RESET POIN pada database elastic

GET /promo-point-exchange-transaction/\_search { \"query\": { \"bool\":
{ \"must\": \[ { \"exists\": { \"field\": \"idtransaksi\" } }, {
\"match\": { \"idcustomer\": \"BZTJS1678156967CFC\" } }, { \"bool\": {
\"must_not\": { \"exists\": { \"field\": \"deleted_at\" } } } } \] } },
\"sort\": \[ { \"created_at\": { \"order\": \"asc\" } } \] }

GET /promo-point-resume/\_search { \"query\": { \"match\": {
\"idcustomer\": \"BZTJS1678156967CFC\" } } }

RESET POIN POST /promo-point-resume/\_update_by_query { \"script\": {
\"source\": \"ctx.\_source.saldo_point = 0;\", \"lang\": \"painless\" },
\"query\": { \"term\": { \"idcustomer.keyword\": \"xxxx\" } } }

id settlement QHY231211100640879 = STEDC17042764443740
BPR231211100253781 = STEDC17042764443740 IRA231211094959264 =
STEDC17042764443740 AOI231211094618291 = STEDC17042764443740
