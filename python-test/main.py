import json
import time
import requests
import pandas as pd
def get_trace_info(feedback):
    trace_data = []
    for key, value in feedback['Requests'].items():
        for i in range(value["Count"]):
            trace_data.append({
                "RequestType":key, 
                "Id":value["TraceIds"][i],
            })

    trace_data = {"data": trace_data}
    start = time.time()
    res = requests.post("http://136.159.209.204:8657/start", json=trace_data)
    end = time.time()
    result = res.json()
    print(len(result))
    return result
def handle_traces(trace_info):
    trace_df = {}
    for r in ['login', 'editbook_getbook', 'editbook', 'getbook']:
        temp = {}
        for trace in trace_info:
            if trace['request_type'] in r:
                for key,value in trace.items():
                    if key not in temp:
                        temp[key] = []
                    
                    if type(value) == float or type(value) == int:
                        # print(value, value / 1000)
                        value = value / 1000
                        temp[key].append(value)
                    else:
                        temp[key].append(value)
        trace_df[r] = pd.DataFrame(temp)
    return trace_df

with open('./input.json', 'r') as f:
    feedback = json.load(f)
trace_info = get_trace_info(feedback)
trace_dfs = handle_traces(trace_info)
if 'login' in trace_dfs.keys():
    # for key in trace_dfs['login'].columns:
    #     print(key)
    print(trace_dfs['login'].describe())