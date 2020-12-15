# shadowsocks-windows-dynamic
[![](https://img.shields.io/badge/Telegram-Group-blue)](https://t.me/aioCloud)
[![](https://img.shields.io/badge/Telegram-Channel-green)](https://t.me/aioCloud_channel) 

```csharp
namespace Shadowsocks
{
    public static class NativeMethods
    {
        [DllImport("shadowsocks.bin", CallingConvention = CallingConvention.Cdecl)]
        public static extern bool ServerInfo(byte[] client, byte[] remote, byte[] passwd, byte[] method);

        [DllImport("shadowsocks.bin", CallingConvention = CallingConvention.Cdecl)]
        public static extern bool Create();

        [DllImport("shadowsocks.bin", CallingConvention = CallingConvention.Cdecl)]
        public static extern void Delete();
    }

    public class Shadowsocks
    {
        public static void Create()
        {
            var client = Encoding.UTF8.GetBytes("0.0.0.0:1080");
            var remote = Encoding.UTF8.GetBytes("1.1.1.1:80");
            var passwd = Encoding.UTF8.GetBytes("114514");
            var method = Encoding.UTF8.GetBytes("chacha20-ietf");

            if (!NativeMethods.ServerInfo(client, remote, passwd, method))
            {
                Console.WriteLine("!NativeMethods.ServerInfo");
                Console.ReadLine();
                return;
            }

            if (!NativeMethods.Create())
            {
                Console.WriteLine("!NativeMethods.Create");
                Console.ReadLine();
                return;
            }
        }

        public static void Main(string[] args)
        {
            Create();
            Console.WriteLine("STARTED");
            Console.ReadLine();

            NativeMethods.Delete();
            Console.WriteLine("STOPPED");
            Console.ReadLine();
        }
    }
}
```
