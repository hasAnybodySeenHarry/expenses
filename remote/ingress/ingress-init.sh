helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update

helm install my-ingress-controller ingress-nginx/ingress-nginx