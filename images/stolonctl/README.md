`stolonctl` - stolon's CLI utility. Basic container behaviour - do noting and wait.
Can be added to main pod, which will allow to execute it by operations on demand
using `kubectl exec -it <stolonctl-pod-name> /stolonctl`. Useful for getting stolon's
internal cluster information.
