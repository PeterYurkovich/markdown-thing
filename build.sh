rm -rf markdown
gh repo clone peteryurkovich/obsidian-notes markdown

sudo docker buildx build . -t markdown-thing:latest
