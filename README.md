Go challenge
============

Skriv en http-server i Go ved hjelp av [standardbiblioteket](https://pkg.go.dev/std).

Serveren skal motta og validere JSON på formen

```json
{
    "seeds": [
        "abc",
        "def",
        "xyz"
    ]
}
```

(på f.eks. `/hash`). `"seeds"` og hver seed-verdi kan ha vilkårlig lengde. Handler skal deretter prosessere hver seed-verdi parallelt i en _goroutine_ som kalkulerer en [SHA256 hash](https://pkg.go.dev/crypto/sha256) av seed-strengen. Prosesseringsrutinen skal legge resultatet på en _channel_.

Bruk for eksempel en [WaitGroup](https://pkg.go.dev/sync#WaitGroup) for prosesseringsrutinen. Når alle seeds er prosessert, skal resultatet returneres som JSON på formen

```json
{
    "hashes": [
        "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad",
        "cb8379ac2098aa165029e3938a51da0bcecfc008fd6795f401178647f96c5b34",
        "3608bca1e44ea6c4d268eb6db02260269892c0b42b86bbf1e77a6fa16c3c9282"
    ]
}
```

der hash-verdiene bruker heksadesimal notasjon og ligger **samme rekkefølge** som i "seeds".

Dersom input er ugyldig, skal serveren returnere en relevant feilkode og feilmelding til klienten.

Ta i bruk relevante headere, som `Content-Type`.

Serveren skal logge tidsstemplede hendelser til `stdout` og feil til `stderr`.

Koden skal ha [tester](https://pkg.go.dev/testing), og alle tester skal passere. Koden skal selvsagt også kompilere og være kjørbar.

Du skal bruke Git for versjonkontroll, og vi foretrekker om du gjør flere commits i løpet av utviklingen, så vi kan få et inntrykk av hvordan du arbeider. Skriv gode commit-meldinger på engelsk. Husk at det ikke er noen skam å snu, eller å endre koden underveis.

Når du er klar til å levere, kan du gjøre det via [pull request](https://github.com/protojour/go-challenge/pulls). Du bør altså starte med en fork.
