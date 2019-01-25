# jscale

地震の加速度記録から計測震度と震度階を求めます。

## Install

``` go get github.com/takatoh/jscale```

## Usage

地震の加速度記録は CSV 形式で、3成分必要です。
たとえば次のように：

```
Time,NS,EW,UD
0.00,-0.03,0.01,-0.34
0.01,0.18,0.42,0.05
0.02,-0.17,-0.25,0.34
0.03,-0.55,0.41,0
0.04,0.31,0.76,0.62
0.05,-0.72,0.1,0.37
0.06,-0.04,0,-0.3
0.07,-0.13,0.46,-0.09
0.08,-0.69,-0.22,0.36
（後略）
```
加速度記録のファイルを example.csv とすると、次のように実行します。

``` jscale example.csv```

または、気象庁や K-NET の強震記録を利用することもできます。
気象庁の強震記録を利用する場合は、-jma オプションを付けて実行します。

``` jscale -jma example.txt```

K-NET の強震記録は成分ごとにファイルが分かれていて、それぞれ example.NS, example.EW, 
example.UD とすると、拡張子を除いた部分を -knet オプションとともに指定します。

``` jscale -knet example```

固定長フォーマットの加速度記録にも対応しました。固定長フォーマットはパラメータが多いので、
TOML 形式の入力ファイルを作成します。たとえば次のように：

```toml
[[wave]]
name   = "NS"
file   = "exmaple.dat"
format = "10F8.2"
dt     = 0.01
ndata  = 12000
skip   = 2

[[wave]]
name   = "EW"
file   = "exmaple.dat"
format = "10F8.2"
dt     = 0.01
ndata  = 12000
skip   = 1204

[[wave]]
name   = "UD"
file   = "exmaple.dat"
format = "10F8.2"
dt     = 0.01
ndata  = 12000
skip   = 2406
```

この例では、example.dat が加速度記録のファイルです。dt と ndata はそれぞれ時刻刻みとデータ数です。
これらが同一の加速度記録が3成分必要です。

実行は次のように -fixed オプションに続けて入力ファイルを指定します。

``` jscale -fixed input.toml```

## License

MIT License
