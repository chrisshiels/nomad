# 'Makefile'.
# Chris Shiels.


CN = vm1
SUBJ = /C=GB/ST=London/L=London/O=Mecachis/OU=Mecachis/CN=$(CN)


$(CN).crt:	$(CN).key
	openssl req \
	-x509 \
	-new  \
	-key $(CN).key \
	-nodes \
	-sha256 \
	-subj $(SUBJ) \
	-days 365 \
	-out ./$(CN).crt


$(CN).key:
	openssl genrsa -out $(CN).key 4096


clean:
	rm -f $(CN).crt $(CN).key
