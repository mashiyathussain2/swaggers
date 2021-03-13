package app

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vasupal1996/goerror"
)

const (
	imgJpeg string = "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/7QBmUGhvdG9zaG9wIDMuMAA4QklNBAQAAAAAAEkcAVoAAxslRxwCAAACAAAcAnQANcKpIGhpcHN0dWZmIC0gaHR0cDovL3d3dy5yZWRidWJibGUuY29tL3Blb3BsZS9oaXBzdHVmAP/bAEMABgQFBQUEBgUFBQcGBgcJDwoJCAgJEw0OCw8WExcXFhMVFRgbIx4YGiEaFRUeKR8hJCUnKCcYHSsuKyYuIyYnJv/bAEMBBgcHCQgJEgoKEiYZFRkmJiYmJiYmJiYmJiYmJiYmJiYmJiYmJiYmJiYmJiYmJiYmJiYmJiYmJiYmJiYmJiYmJv/CABEIAIAAgAMBIgACEQEDEQH/xAAbAAEAAwEBAQEAAAAAAAAAAAAAAwQGBQcBAv/EABkBAAMBAQEAAAAAAAAAAAAAAAADBAUCAf/aAAwDAQACEAMQAAAB9QASU+bje5vQ+JyNP7zw4u/xXRRVujreW5Tq/c1y3UdbzfZ8t6IXSAcHLel4h+dVpdvr8s7mTt+aee3d9gPQOW97EWvtGfQ2Ge1yqowq0Akwe8wr8+2KcmvN+qUul0qcv695CnPj2eJ20uxGE3AEmJ1+PdDYFWOAAABX2+G10uxOE3AGO/G4xPcc4txQAAAKlqnvI9eMLvAJMPuMK6C2KsgAAAClvMJu5tWMI0QCTE7bKtkiFeIAAABHs81pY9yMLqAJOf0Mv4cu7PK3Ppv3HTD9fJj2KjoKk13cucDvqvjHp//EACYQAAICAQMDBAMBAAAAAAAAAAMEAQIFABIgEBMwERQVNAYhMSL/2gAIAQEAAQUC6kaWFdhoK9LtOF1NL217dfXtl9b1aSBxiI+TU2ifXvfjkjX9fepjYXXmt6Y1SNFxq0wObaEMZtDikUeTgurRJhoH98pjSXIpwyQ71IoFaX2pYgH4d3fTRpj5DMKtMY7BL3WRtMVhX69O5DKgPbr8Qzuid0XpkRehchMxk8mdYmEyLd7xkCehiHZ6AnZkOSX0+h43DYaHZch6D0CLVD0mdp+Sn6W6l70pELdg/Uv95B8JuRSiDQV6Xv4GLRSgGAMRw7gxuhm1r+Bv1he5JM7xF/fA79Lkt+w+B36XJT63gc/anIUbZ8ExvNydp2XvBjKdw/J+ASqMhKxyt3WKJ3Fdbllu7Q/+nNfGdusrvV1ta1tb1VV6+qY0U6uK6k4kZ4r1/8QAJxEAAQMDAgUFAQAAAAAAAAAAAQIDBAAQEhEhEyAyM0EiIzFxgWH/2gAIAQMBAT8Bpx/A4ga0FOPDVGwp9t1G5VqKZjrIyKtKW44xsremXislKvkXlpCRmPO1RnkFGOtPHijhooOJTss6GpTyV+lNREpwy83ndItDfDYKTUx/iHEfAtC7f7eb2/3lhds/d5qvaA/vLBV6CLzegcsHpN5acm/rliowb5H4Q11bNGO6PFCM6fFMQwDq5f8A/8QAIREAAQQBBAMBAAAAAAAAAAAAAQACAxEQICEiMRITMkH/2gAIAQIBAT8BTWXuVQb2mOafxOeOqQaH9Jza3GYjeykabtN4myqvpRMI3Kku6zD3ieHzIKghDN8S/WYfrTN9ZgZzLtM7bcDmHvTN3mI07TIbdk7BRTkjmvY1exqlmdXBDH//xAA9EAACAQEDBgkKBgMBAAAAAAABAgMRAAQSICExQVFhEyIwMlJxgZHBEBQjQnOhorGy0TM0U2JykgXC8NL/2gAIAQEABj8C8uCS8RI3RZwDYNI3O5qjOW6hbiKt2Tfx2+w99qyz3iTrkwj4aW/AjP8AJa/O35aLsQWIS9GIj9OY5uzRbEkqX2PsDd4zfKw47Y/0gpLjrFhGccTtoEiFa5Ud2jbBwgLOw0hRs77NcYIC8ujCqihO8+NlhgCcO4JZ8OZF3DZsFvSRCdulNxjasC+bPqaLN3jQbMsgCyIcLgbbPLeXCXSNsNGNA5113aqWAjACasOixmhAS8gZm6W47rJeLvxZcPFxa/2tZ/OIcJDYHTVYcIcToWQttoaVyY70ql1ClHCipptteL3BQ1AGIGoxHOfC3+Q81xcPwMeHDppiatLFlneSJlONCuZHrmz66ivkvBGgBAevP4EWu4u6NLwEkiyImc1rpp/2m3BtjCYiY1k5yrvsSxoBrNkNKYqtTrJPjaSCKIyM/pF1KK5jU9Y99hFixnOWbaTnOVJJ05Xb4qfICyyxECRNFdBGsG3popom9mWHetj5vE3tJhgUd+c24G7a+OZnTPITrs8rFVOZS7LRJNx2HfbjXJ8X7ZFp87YZsKRa4lNcXWfDyXc9MPH4/wCuXD/Hy83HxlJXpCuizz+aveI461LoAPitwF4gdeL+GExhh2WRX0gba08t1bZOPeCPHLjHRqO4nIluSRijuWWQtmAJxZ9emyytCYQkZWjMDUkjZqzZEPt4/qy5Rsmk+s8jF7aP6xlY5pFjXaxpa8tG4dOHNCpqNA5FXbQssZP9xYmGVXppodGTOb1Thw1E9ZiurCuynvra8sylCZjxTqzLyLMNKlT3MLXcRRkShuMWUhkXXqzjV3ZUx2zyfVTw5Gf2ZOXi6TOfjPI3j2TfLLTt+o8jMvSXD35suWPoTOPfX5HkbvF0pQexeN4DLxepeB8Y+4+nkZLz6i+ij39I9+bsy3F4bCm0aa6qb7Kt6Xgmbms2YP8AY7ss+aiqetJXDi3Idu/RZDAMMYFAtKYaasuGSuFKYRKdERPrHszCxjgpwB58rCo6l29ejrsBd71ItNUnHH/dRtnhik/hJT3EeNvyM39k/wDVvyMva6DxtzYIRtLFz4WreXa8/tbMn9fvW3B0ka6sKK0YLNFu6tmy0ks4wmSmbaaUxbq7Mj//xAAnEAEAAQMCBQUBAQEAAAAAAAABEQAhMUFRIGFxgZEwobHB8BDh8f/aAAgBAQABPyH+GK/5zsBavCbCz9gu0Yu5D8L71M5TcTtBWut3B/Kh2CXkvgrXHgjrnJ8ir/tvcfS5eSdVDtI7KR2Ek6451O6ETx2FsvKZpxwhTpBxBEbFQToTrFBQBTy75OqHWhWLIUy+xJjUbsw0VOs715sdACkA47CP5MlLlz4DVJyRE61otQE2E9r1EZm1C4hbZDlFGbFaPnr2yU2QpdIt+YT4e8tG0b4SQ21G/ZoylNt93EB4YgciBIhNYS4Xh5VCF3lFxO3uak4E37iJxyoSQSwGNiykMcp0/nwYFiaIfYCIOZBmRn/VEQdFczBG05e9DUNKkAVdYkFaBh4FTliH8EhctK4FYusoROUjS624TFdYK5SFAT9fcXhMF9EHki54HggT9ijbGu0e8Q8AcypzoC6LgmSI0saWtFZGdEGjltkdEbUGHUD5IfajYd+Xsy3IX1UtTdlrnungPy8044TFXruX3f6sBEB9Ar3B9Vc6IIkYQymbWKVZKo0DsYxcyGanIMRkGguqEE8v7yBNM04TFH8pF9cEp8AxOQNVxYiIvUz2GBCxJs16zjgzdvhH3WnCYqwwxfc+/ScccOiT458tNyiyhK56r6LWQ7RNinCb5mes078N829WcuCJLEpm5irpMJFsrwppv6LDZAgxMj6q0Sb2rhbFJyiVGOExSlf+BeiZu3hE/VOHhMUov9Kvo/u76cPCYrH2j6IYTMDqj7U4eExUhEMFyWfZ9GH9b8j3DvTjhMVyiidhx6K1pFTq+Ad3enHCYpszxtZaN7oiNadeqAGwxbf8Z0bMPEiEHkTX/It5uXNISTBM0SIinHCYpgTOIF4mTW9MSwxNWR4ormw2fNZ4GLyDP9mE7ByrEnuzaDf7T02v1GHUBKGiXYD71GKM2Q7LNDHxcQYldy+htFNxa5EJFI1QOSKcf3//2gAMAwEAAgADAAAAEABrolbIAAFTuq8aQAEP7LP4QAAf/wD/APmAAF//AP8A9MAAL/8A/wD4wAEP/wD/APBAABW5wxJA/8QAJBEBAAEBBgcBAAAAAAAAAAAAAREAECExYXHBIEFRgZGx8KH/2gAIAQMBAT8QrGWCXlBrUoM7FdDp+0bO1H1QlEuQ0oIBwW6HOKBKOhhawSG8z57RRFAT66idRXF5Bm9ctaFCAaeKNEkL5q4OK57fTaPPtY44Dff4qHieYadbFL1ei0SXoPTw/RkW9uvTwsZ/a39mzwi9ztVU5p23nhUZxb/u0WgLDUr2Ts/a1ij9+qwvzgou85GHf7zTjZ//xAAiEQEAAgICAQQDAAAAAAAAAAABABEQITFBIFFhcdGBofD/2gAIAQIBAT8QlGyiK63ZqyjDtAzdaJDIljlFbg3BMGmem4uybJYwfsjPJgNzZqEou3B/TKr8PHg+PvIHcT68RV5M8vitDPId+NtXWUok1V/fEF7iPcuaRKW4/8QAJRABAQACAQMEAwEBAQAAAAAAAREAITFBUXEgMGGhEIGRsfDh/9oACAEBAAE/EPxw+MtVYi1eNBxsCtHAu0WbYQNqG8XU2BXO6Avxjuzar8OT+pyz3UQv3TB1ZwO/Q5TA9pB2BeMCutfmgiEXvO8ElihacpvOBURCKqvjWkbxC0dZyePTY12xiJ09c6SBA8lAj8tVRuoRKzYdpmASuoQwwpvKSaTXGdwMfhHQM4TFPT4NPRG1wjEAcGtSACu99thDscfWFMuQEGkVGQIwtyYgOmNJ4yn4UYDifSTdW6EyJJps7SE2EG3eIRg3HFQWlKCOqgIjIyPn0IxfIHyvpc3Nwex2ZQNKDo33SgGfCtF+R1xigwb0JLqol1Fjcw4gxUCWIjBJDt+CDaKHUO8/YMsa3+EezRGaA8LN3Xl7iLVdg4kYiCFZGAtqrwBu4lhS0i4x0mp0ytGITYp0u0FSm4byUCuSLgsKYQrPTw+MIRs2JcKobiuyToV3MChBjgYj4Cs9VzLxUe44zmKDBkq9S5gjd5sEDaG00Ny6MIBQMqkTHdhmAdQZ0GL5/Tuf6X5L8Y3f+RleGB42LgopKJVavfICX+g/nfazk8enh/WJW5U8qP2v5ubOFTnsFRp08ucZ1MEDI2QUKUlykTu2UGTRBOKp1w/IR3gfnRNqprz+UA02PjDPPx6eH9ZRSIfz6FNge6ABVRQUSSImJYAaHtJDRIrilfyoD/1Zjz8enh8YNRAp2Ef0fZ5Y3k8ektRDfl4KBdYBLFlKQafZDDDSSKbAV10MpY4DQ2HM9OgekaDZGrFyRGmPSOU7/j4IKxKgpZyeyrg7hIcFgLpVDvMate21sQlC07dj0Ck+MTnDb4P6PstBckPkwH8H08P6x0Ntf/o59q79D08P6wRfNvkF+z2X6Q33En9OfSfTw+MXxVR9lNuki5llU8fzucvj08XjECYI6FNvkSfpynsAqAKugOuW9ghwgp9xGfm4Gcnj08XjJE5CpGVBRKIVAA8YCOjyGw0pCo8u2xgKAiciRPSC2HBV6B3exjwUQKAxgiilRNGC+YppzVO2pLhJnJ49PF4yTRniImggyIsDfeObcv2W98D2lcKFMkAx/g9BYJUWwD3eoPFYdx/Ff1P+ZZhfkb+p+sSARv8A45DfIQxI0N8Gf2DvhzclIomqLyBTRKEOm9bI83atgWKhyePz/9k="
)

func TestIMG_DecodeBase64StrToIMG(t *testing.T) {

	type args struct {
		b64Str string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "[Error] Empty base64 string",
			args: args{
				b64Str: "",
			},
			wantErr: true,
			err:     goerror.New("invalid base64 image string [format should be`data:image/(jpeg/png);base64,/9j/4AAQSkZJRgABAQE....`]", nil),
		},
		{
			name: "[Error] Only meta png data",
			args: args{
				b64Str: "data:image/png;base64",
			},
			wantErr: true,
			err:     goerror.New("invalid base64 image string [format should be`data:image/(jpeg/png);base64,/9j/4AAQSkZJRgABAQE....`]", nil),
		},
		{
			name: "[Error] Only meta jpeg data",
			args: args{
				b64Str: "data:image/jpeg;base64",
			},
			wantErr: true,
			err:     goerror.New("invalid base64 image string [format should be`data:image/(jpeg/png);base64,/9j/4AAQSkZJRgABAQE....`]", nil),
		},
		{
			name: "[Error] Only meta png data",
			args: args{
				b64Str: "data:image/png;base64",
			},
			wantErr: true,
			err:     goerror.New("invalid base64 image string [format should be`data:image/(jpeg/png);base64,/9j/4AAQSkZJRgABAQE....`]", nil),
		},
		{
			name: "[Error] empty image jpeg data",
			args: args{
				b64Str: "data:image/jpeg;base64,",
			},
			wantErr: true,
			err:     goerror.New("invalid base64 image string [format should be`data:image/(jpeg/png);base64,/9j/4AAQSkZJRgABAQE....`]", nil),
		},
		{
			name: "[Ok]",
			args: args{
				b64Str: imgJpeg,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &IMG{}
			err := i.DecodeBase64StrToIMG(tt.args.b64Str)
			if (err != nil) != tt.wantErr {
				t.Errorf("IMG.DecodeBase64StrToIMG() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				assert.Equal(t, tt.err, err)
			}
			if !tt.wantErr {
				assert.Equal(t, "image/jpeg", i.Type)
				assert.Equal(t, 128, i.Conf.Height)
				assert.Equal(t, 128, i.Conf.Width)
			}
		})
	}
}

func TestIMG_Resize(t *testing.T) {
	type fields struct {
		Type string
		Conf *image.Config
		Img  *image.Image
	}
	type args struct {
		width  uint
		height uint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *image.Image
	}{
		{
			name: "Passing original image height and width",
			args: args{
				height: 128,
				width:  128,
			},
		},
		{
			name: "Downscaling image",
			args: args{
				height: 100,
				width:  100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &IMG{}
			i.DecodeBase64StrToIMG(imgJpeg)
			got := i.Resize(tt.args.width, tt.args.height)
			assert.NotNil(t, got)
		})
	}
}

func TestGenerateOTP(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "[Ok] Passing 0",
			args: args{
				length: 0,
			},
			wantErr: false,
			want:    "",
		},
		{
			name: "[Ok] Passing 10",
			args: args{
				length: 10,
			},
			wantErr: false,
			want:    "0987654323",
		},
		{
			name: "[Ok] Passing 1",
			args: args{
				length: 1,
			},
			wantErr: false,
			want:    "0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateOTP(tt.args.length)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateOTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Len(t, got, len(tt.want))
			}
		})
	}
}
