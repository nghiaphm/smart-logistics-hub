import NextAuth from "next-auth";
import KeycloakProvider from "next-auth/providers/keycloak";

console.log("👉 SOI CLIENT ID:", process.env.KEYCLOAK_CLIENT_ID);
console.log("👉 SOI SECRET:", process.env.KEYCLOAK_CLIENT_SECRET);
const handler = NextAuth({
  providers: [
    KeycloakProvider({
      clientId: process.env.KEYCLOAK_CLIENT_ID!,
      clientSecret: process.env.KEYCLOAK_CLIENT_SECRET!,
      issuer: process.env.KEYCLOAK_ISSUER,
    }),
  ],
  // Tùy chọn: Thêm các callbacks nếu sau này bạn muốn lấy Token để gửi sang Golang
  callbacks: {
    async jwt({ token, account }) {
      if (account) {
        token.accessToken = account.access_token;
      }
      return token;
    },
    async session({ session, token }) {
      // Ép kiểu an toàn để đẩy accessToken vào session
      (session as any).accessToken = token.accessToken;
      return session;
    },
  },
});

export { handler as GET, handler as POST };