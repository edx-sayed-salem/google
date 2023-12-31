package play

import (
   "154.pages.dev/protobuf"
   "errors"
   "io"
   "net/http"
)

type Details struct {
   m protobuf.Message
}

// play.google.com/store/apps/details?id=com.google.android.youtube
func (d Details) Downloads() (uint64, bool) {
   d.m.Message(13)
   d.m.Message(1)
   return d.m.Varint(70)
}

func (d Details) Files() []uint64 {
   var files []uint64
   d.m.Message(13)
   d.m.Message(1)
   for _, f := range d.m {
      if f.Number == 17 {
         if m, ok := f.Message(); ok {
            if file, ok := m.Varint(1); ok {
               files = append(files, file)
            }
         }
      }
   }
   return files
}

// play.google.com/store/apps/details?id=com.google.android.youtube
func (d Details) Name() (string, bool) {
   return d.m.String(5)
}

// play.google.com/store/apps/details?id=com.google.android.youtube
func (d Details) Offered_By() (string, bool) {
   return d.m.String(6)
}

// play.google.com/store/apps/details?id=com.google.android.youtube
func (d Details) Price() (float64, bool) {
   d.m.Message(8)
   if v, ok := d.m.Varint(1); ok {
      return float64(v) / 1_000_000, true
   }
   return 0, false
}

// play.google.com/store/apps/details?id=com.google.android.youtube
func (d Details) Price_Currency() (string, bool) {
   d.m.Message(8)
   return d.m.String(2)
}

// play.google.com/store/apps/details?id=com.google.android.youtube
func (d Details) Requires() (string, bool) {
   d.m.Message(13)
   d.m.Message(1)
   d.m.Message(82)
   d.m.Message(1)
   return d.m.String(1)
}

func (d Details) Size() (uint64, bool) {
   d.m.Message(13)
   d.m.Message(1)
   return d.m.Varint(9)
}

// play.google.com/store/apps/details?id=com.google.android.youtube
func (d Details) Updated_On() (string, bool) {
   d.m.Message(13)
   d.m.Message(1)
   return d.m.String(16)
}

// developer.android.com/guide/topics/manifest/manifest-element
func (d Details) Version_Code() (uint64, bool) {
   d.m.Message(13)
   d.m.Message(1)
   return d.m.Varint(3)
}

// developer.android.com/guide/topics/manifest/manifest-element
func (d Details) Version_Name() (string, bool) {
   d.m.Message(13)
   d.m.Message(1)
   return d.m.String(4)
}

func (h Header) Details(doc string) (*Details, error) {
   req, err := http.NewRequest("GET", "/fdfe/details?doc=" + doc, nil)
   if err != nil {
      return nil, err
   }
   req.URL.Scheme = "https"
   req.Header.Set(h.Device_ID())
   req.Header.Set(h.Agent())
   req.Header.Set(h.Authorization())
   req.URL.Host = "android.clients.google.com"
   res, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return nil, errors.New(res.Status)
   }
   mes, err := func() (protobuf.Message, error) {
      b, err := io.ReadAll(res.Body)
      if err != nil {
         return nil, err
      }
      return protobuf.Consume(b)
   }()
   if err != nil {
      return nil, err
   }
   mes.Message(1)
   mes.Message(2)
   mes.Message(4)
   return &Details{mes}, nil
}
