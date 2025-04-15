#include <libcryptage.h>
#include <stdio.h> // for printf
#include <stdlib.h> // for free

#define PUBKEY "age1mcka6j3umwgqklmcxsph4de0g7ar7wqnelxrgzytj4rfwc8pldrs3nu3qx"
#define PRIVKEY "AGE-SECRET-KEY-1MKEHFSSMELSPTCZUYYKMZML5J6H7N7AJTRNMXQXQ60D8MHFLM37S0DR2KJ"

int main(){
    char* got = age_encrypt_armor(PUBKEY, "hello, libcryptage", 0);
    if (got == NULL) {
        printf("%s\n", ageerr());
        return 1;
    }
    printf("encrypted:\n%s\n", got);
    char* dec = age_decrypt_armor(PRIVKEY, got, 0);
    if (dec == NULL) {
        printf("%s\n", ageerr());
        return 1;
    }
    printf("decrypted: %s\n", dec);

    // free returned strings
    free(got);
    free(dec);
    // but dont free errors returned by ageerr() !!!
    return 0;
}

