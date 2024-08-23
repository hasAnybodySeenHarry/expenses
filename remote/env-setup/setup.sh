helm repo add argo https://argoproj.github.io/argo-helm
helm repo update

kubectl create namespace argocd

helm install argocd argo/argo-cd --namespace argocd

helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update

helm install my-ingress-controller ingress-nginx/ingress-nginx --set controller.allowSnippetAnnotations=true --set allowSnippetAnnotations=true

kubectl apply -f https://raw.githubusercontent.com/hasAnybodySeenHarry/debt-tracker/main/remote/postgres/postgres.yaml
kubectl apply -f https://raw.githubusercontent.com/hasAnybodySeenHarry/debt-tracker/main/remote/postgres/setup-job.yaml

kubectl apply -f https://raw.githubusercontent.com/hasAnybodySeenHarry/debt-tracker/main/remote/amqp/amqp.yaml

kubectl apply -f https://raw.githubusercontent.com/hasAnybodySeenHarry/debt-tracker/main/remote/redis/redis.yaml

kubectl apply -f https://raw.githubusercontent.com/hasAnybodySeenHarry/debt-tracker/main/remote/ratelimit/ratelimit.yaml

kubectl apply -f https://raw.githubusercontent.com/hasAnybodySeenHarry/debt-tracker/main/remote/pub-sub/zookeeper.yaml
kubectl apply -f https://raw.githubusercontent.com/hasAnybodySeenHarry/debt-tracker/main/remote/pub-sub/kafka.yaml

kubectl apply -f https://raw.githubusercontent.com/hasAnybodySeenHarry/debt-tracker/main/remote/argo/repo.yaml
kubectl apply -f https://raw.githubusercontent.com/hasAnybodySeenHarry/debt-tracker/main/remote/argo/app.yaml

kubectl apply -f https://raw.githubusercontent.com/hasAnybodySeenHarry/debt-tracker/main/remote/ingress/ingress.yaml
