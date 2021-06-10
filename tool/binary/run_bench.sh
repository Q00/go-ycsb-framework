rm -rf ./data 
rm -rf ./logs

for db in cockroach mariadb tidb tikv raw cassandra scylla
do
    ./bench.sh load ${db}
    ./bench.sh run ${db}
done

./clear.sh