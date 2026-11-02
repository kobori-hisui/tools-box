// Package process is reserved for supervised child processes (yt-dlp, ffmpeg,
// converters): start with a context, cap output, kill the whole process group
// on cancel, and never build a shell command line from extension input.
package process
