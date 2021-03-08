package foundation.icon.test.cases;

import foundation.icon.icx.IconService;
import foundation.icon.icx.KeyWallet;
import foundation.icon.icx.transport.http.HttpProvider;
import foundation.icon.test.common.Constants;
import foundation.icon.test.common.Env;
import foundation.icon.test.common.TestBase;
import foundation.icon.test.common.TransactionHandler;
import foundation.icon.test.score.ChainScore;
import foundation.icon.test.score.GovScore;
import foundation.icon.test.score.IISSScore;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import foundation.icon.icx.data.TransactionResult;
import foundation.icon.icx.data.Address;

import java.math.BigInteger;

import static foundation.icon.test.common.Env.LOG;

@Tag(Constants.TAG_JAVA_INTEGRATION)
public class IISSTest extends TestBase {
    private static final String SCORE_STATUS_PENDING = "pending";
    private static final String SCORE_STATUS_ACTIVE = "active";
    private static final String SCORE_STATUS_REJECTED = "rejected";

    private static TransactionHandler txHandler;
    private static ChainScore chainScore;
    private static GovScore govScore;
    private static GovScore.Fee fee;
    private static KeyWallet[] testWallets;
    private static final int testWalletNum = 3;
    private static KeyWallet governorWallet;
    private static IISSScore score;
    private Address score_address;

    enum TargetScore {
        TO_CHAINSCORE(Constants.CHAINSCORE_ADDRESS),
        TO_GOVSCORE(Constants.GOV_ADDRESS);

        Address addr;
        TargetScore(Address addr) {
            this.addr = addr;
        }
    }

    @BeforeAll
    public static void init() throws Exception {
        Env.Node node = Env.nodes[0];
        Env.Chain chain = node.channels[0].chain;
        IconService iconService = new IconService(new HttpProvider(node.channels[0].getAPIUrl(Env.testApiVer)));
        txHandler = new TransactionHandler(iconService, chain);
        chainScore = new ChainScore(txHandler);
        govScore = new GovScore(txHandler);
        fee = govScore.getFee();
        governorWallet = chain.governorWallet;

        BigInteger ownerBalance = ICX.multiply(new BigInteger("30000")); // deploy(100) + deposit(5000)
        testWallets = new KeyWallet[testWalletNum];
        for (int i = 0; i < testWalletNum; i++) {
            testWallets[i] = KeyWallet.create();
            txHandler.transfer(testWallets[i].getAddress(), ownerBalance);
        }
        score = IISSScore.install(txHandler, testWallets[0]);
        Address score_address = score.getAddress();
    }

    @Test
    public void registerPRepAndSetBonderListByRPC() throws Exception {
        LOG.infoEntering("registerPRepAndSetBonderListByRPC");
        String name = "ABC";
        String email = "abc@example.com";
        String country = "KOR";
        String city = "Seoul";
        String website = "https://abc.example.com/";
        String details = "https://abc.example.com/details/";
        String p2pEndpoint = "123.45.67.89:7100";
        String nodeAddress = testWallets[1].getAddress().toString();
        TransactionResult result = chainScore.registerPRep(testWallets[1], name, email, country, city, website, details, p2pEndpoint, nodeAddress, ICX.multiply(new BigInteger("2000")));
        assertSuccess(result);

        String [] arrayAddress = {score.getAddress().toString()};
        result = chainScore.setBonderList(testWallets[1], arrayAddress);
        assertSuccess(result);
        LOG.infoExiting();
    }

    @Test
    public void setStake() throws Exception {
        LOG.infoEntering("setStake");
        String val = "1000";
        TransactionResult result = score.setStake(testWallets[0], val);
        assertSuccess(result);
        LOG.infoExiting();
    }

    @Test
    public void setDelegation() throws Exception {
        LOG.infoEntering("setDelegation");
        String val = "300";
        TransactionResult result = score.setDelegation(testWallets[0], testWallets[1].getAddress(), val);
        assertSuccess(result);
        LOG.infoExiting();
    }

    @Test
    public void setBond() throws Exception {
        LOG.infoEntering("setBond");
        String val = "300";
        TransactionResult result = score.setBond(testWallets[0], testWallets[1].getAddress(), val);
        assertSuccess(result);
        LOG.infoExiting();
    }

    @Test
    public void getStake() throws Exception {
        LOG.infoEntering("getStake");
        Object obj = score.getStake(testWallets[0], score.getAddress());
        System.out.println(obj);
        LOG.infoExiting();
    }

    @Test
    public void getBalanced() throws Exception {
        LOG.infoEntering("getBalanced");
        TransactionResult result = score.getBalance(testWallets[0]);
        assertSuccess(result);
        LOG.infoExiting();
    }

    @Test
    public void getPrep() throws Exception {
        LOG.infoEntering("getPrep");
        Object obj = score.getPrep(testWallets[1], testWallets[1].getAddress());
        System.out.println(obj);
        LOG.infoExiting();
    }

/*
    @Test
    public void unregisterPRep() throws Exception {

        TransactionResult result = chainScore.unregisterPRep(testWallets[0]);
        System.out.println(result.getStatus());
        System.out.println(result.toString());
    }*/
}
