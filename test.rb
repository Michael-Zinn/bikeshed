
# You can test a folder full of jpgs like this:
# ls *.jpg | ruby test.rb

ARGF.each do |filename|
	filename.chop!
	puts "[#{filename}]"
	color = `bikeshed #{filename}`
	puts color
	`convert -size 80x120 xc:##{color} #{filename}.bikeshed.png`
end
