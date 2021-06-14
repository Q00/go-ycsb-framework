rm -rf ./data 
rm -rf ./logs

for db in  tikv pg mariadb mongodb cockroach tidb
do
    ./bench.sh load ${db}
    ./bench.sh run ${db}
done

./clear.sh
