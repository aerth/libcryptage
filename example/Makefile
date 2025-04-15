example: example.c
	gcc -o example $^ -I ../ -L.. -lcryptage -static
postinstall: example.c # after installing we shouldnt need the I,L flags
	gcc -o example $^ -lcryptage -static
	@echo "post-install check worked"
clean:
	rm -f example
