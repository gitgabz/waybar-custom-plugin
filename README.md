This is a quick go sample for to create a custom waybar plugin.

It can run either with self updating output or though SIGRTMIN+N signal with a refresh interval.

### Refresh by interrupt signal, restarts the process that has exited.  
The process runs once, outputs and dies. Then it is triggered by interval though signal or directly by SIGRTMIN+8 signal "pkill -RTMIN+8 waybar"
```
"custom/plugin": {
	"format": "{}",
		"max-length": 40,
		"tooltip": false,
		"signal": 8,
		"interval": 12,
		"exec": "~/.local/bin/waybar_plugin -s",
		"on-click": "pkill -RTMIN+8 waybar",
		"return-type": "json"
}
```

### Refresh continuously by Newline json and repeat output at will.
```
"custom/lunch": {
	"format": "{}",
	"max-length": 40,
	"tooltip": false,
	"exec": "~/.local/bin/waybar_plugin",
	"return-type": "json"
}
```