 curl https://www.google-analytics.com/collect \
    -d "tid=UA-150579820-1" \
    -d "t=event" \ 
    -d "ec=testCategory" \ 
    -d "ea=testAction" \ 
    -d "v=1" \ 
    -d "cid=150579820" \
    -o /dev/null
