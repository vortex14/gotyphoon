go get github.com/vortex14/gotyphoon/integrations/python3

archs=(amd64)

for arch in ${archs[@]}
do
        env GOOS=linux GOARCH=${arch} go build -o agent_${arch} agent.go
done


cp agent_amd64 /Users/vortex/ci-agent/agent
sh /Users/vortex/ci-agent/build.sh